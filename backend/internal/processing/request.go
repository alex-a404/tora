package processing

import (
	"context"
	"time"
	pb "tora/backend/gen"
	"tora/backend/internal/telematics"
	"tora/backend/internal/tracker"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserRequest struct {
	FromCoords telematics.Coordinates
	ToCoords   telematics.Coordinates
}

func NewUserRequest(fromCoords, toCoords telematics.Coordinates) UserRequest {
	return UserRequest{
		FromCoords: fromCoords,
		ToCoords:   toCoords,
	}
}

// ProcessReq interacts with python data abstraction layer
func (req UserRequest) ProcessReq(stops []telematics.Stop,
	mgr telematics.Manager,
	trs tracker.TrackingService) UserResponse {

	// determine FROM stop_id and TO stop_id
	fromID, ok := telematics.FindNearestStop(req.FromCoords, stops, 300)
	if !ok {
		return fail("Server-side error (finding closest stop FROM)")
	}

	toID, ok := telematics.FindNearestStop(req.ToCoords, stops, 300)
	if !ok {
		return fail("Server-side error (finding closest stop TO)")
	}

	// TODO find closest bus to FROM
	bus := mgr.GetBuses()[0]

	// connect to gRPC server for python function call
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fail("Server-side error (gRPC)")
	}
	defer conn.Close()

	client := pb.NewOptimizationClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var itineraryList []int32
	for _, stop := range bus.Itinerary {
		itineraryList = append(itineraryList, stop.Id)
	}

	var constraintList []*pb.RouteDependency
	for _, constraint := range bus.Dependencies {
		dep := &pb.RouteDependency{
			From: constraint.From.Id,
			To:   constraint.To.Id,
		}
		constraintList = append(constraintList, dep)
	}

	resp, err := client.RouteOptimize(ctx, &pb.OptimizeRequest{
		Itinerary:   itineraryList,
		Constraints: constraintList,
		Direction:   1,
	})

	return UserResponse{}
}

type UserResponse struct {
	Ok         bool
	FromStop   telematics.Stop
	ToStop     telematics.Stop
	Message    string
	TrackingID int
}

func fail(msg string) UserResponse {
	return UserResponse{
		false,
		telematics.Stop{},
		telematics.Stop{},
		msg,
		0,
	}
}
