package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SimonasKalpokas/todo_server/database"
	pb "github.com/SimonasKalpokas/todo_server/task_service_protos"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedTaskServiceServer
	queries database.Queries
}

func MapReoccurrenceToDbType(r pb.Reoccurrence) (database.Reoccurrence, error) {
	switch r.Type() {
	case pb.Reoccurrence_REOCCURRENCE_NOT_REPEATING.Type():
		return database.ReoccurrenceNotrepeating, nil
	case pb.Reoccurrence_REOCCURRENCE_DAILY.Type():
		return database.ReoccurrenceDaily, nil
	case pb.Reoccurrence_REOCCURRENCE_WEEKLY.Type():
		return database.ReoccurrenceWeekly, nil
	}
	return database.ReoccurrenceDaily, fmt.Errorf("Unexpected reoccurrence value: %s", r.String())
}

func MapTaskTypeToDbType(t pb.TaskType) (database.TaskType, error) {
	switch t.Type() {
	case pb.TaskType_TASK_TYPE_CHECKED.Type():
		return database.TaskTypeChecked, nil
	case pb.TaskType_TASK_TYPE_TIMED.Type():
		return database.TaskTypeTimed, nil
	case pb.TaskType_TASK_TYPE_UNSPECIFIED.Type():
		return database.TaskTypeParent, nil
	}
	return database.TaskTypeChecked, fmt.Errorf("Unexpected task type type value: %s", t.String())
}

func (server *Server) CreateTask(ctx context.Context, createTaskRequest *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	log.Print("Creating task...")
	reoccurrence, err := MapReoccurrenceToDbType(createTaskRequest.Reoccurrence)
	if err != nil {
		return nil, err
	}
	taskType, err := MapTaskTypeToDbType(createTaskRequest.Type)
	if err != nil {
		return nil, err
	}
	createTaskParams := database.CreateTaskParams{
		ID:           ulid.Make().String(),
		Name:         createTaskRequest.Name,
		Description:  createTaskRequest.Description,
		ParentID:     createTaskRequest.ParentId,
		Reoccurrence: reoccurrence,
		Type:         taskType,
	}
	if err := server.queries.CreateTask(ctx, createTaskParams); err != nil {
		return nil, err
	}

	return &pb.CreateTaskResponse{}, nil
}

func (server *Server) CompleteTask(ctx context.Context, completeTaskRequest *pb.CompleteTaskRequest) (*pb.CompleteTaskResponse, error) {
	if err := server.queries.CompleteTask(ctx, completeTaskRequest.TaskId); err != nil {
		return nil, err
	}

	return &pb.CompleteTaskResponse{}, nil
}

func main() {
	// create a new grpc server
	grpcServer := grpc.NewServer()

	ctx := context.Background()

	connStr := os.Getenv("DB_CONNECTION_STRING")
	fmt.Println(connStr)
	conn, err := pgx.Connect(ctx, os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer conn.Close(ctx)

	queries := database.New(conn)
	pb.RegisterTaskServiceServer(grpcServer, &Server{
		queries: *queries,
	})

	// Serve traffic
	log.Print("Serving..")
	wrappedGrpc := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
	)

	httpServer := &http.Server{
		Addr: ":9001",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrappedGrpc.ServeHTTP(w, r)
		}),
	}

	httpServer.ListenAndServe()
}
