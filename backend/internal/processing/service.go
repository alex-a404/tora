package processing

import (
	"context"
	"time"
	pb "tora/backend/gen"
	"tora/backend/internal/telematics"
	"tora/backend/internal/tracker"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	TelematicsMgr   telematics.Manager
	StopManifest    telematics.StopManifest
	TrackingService tracker.TrackingService
}

// ProcessReq interacts with python data abstraction layer
func (s *Service) ProcessReq(req UserRequest) UserResponse {

	// determine FROM stop_id and TO stop_id
	fromStop, ok := telematics.FindNearestStop(req.FromCoords, s.StopManifest.Stops, 300)
	if !ok {
		return fail("Server-side error (finding closest stop FROM)")
	}

	toStop, ok := telematics.FindNearestStop(req.ToCoords, s.StopManifest.Stops, 300)
	if !ok {
		return fail("Server-side error (finding closest stop TO)")
	}

	// TODO find closest bus to FROM
	bus := s.TelematicsMgr.GetBuses()[0]

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

	// add current request to itinerary
	itineraryList = append(itineraryList, fromStop.Id, toStop.Id)

	var constraintList []*pb.RouteDependency
	for _, constraint := range bus.Dependencies {
		dep := &pb.RouteDependency{
			From: constraint.From.Id,
			To:   constraint.To.Id,
		}
		constraintList = append(constraintList, dep)
	}

	// add constraints of current request
	constraintList = append(constraintList, &pb.RouteDependency{
		From: fromStop.Id,
		To:   toStop.Id,
	})

	optimizeResp, err := client.RouteOptimize(ctx, &pb.OptimizeRequest{
		Itinerary:   itineraryList,
		Constraints: constraintList,
		Direction:   1,
	})
	if err != nil {
		return fail("Server-side error (optimizing route gRPC call)")
	}

	var newItinerary []telematics.Stop
	for _, stopId := range optimizeResp.Stops {
		if newStop, ok := s.StopManifest.GetStopFromId(stopId); ok {
			newItinerary = append(newItinerary, *newStop)
		}
	}

	err = bus.NewItinerary(newItinerary)
	if err != nil {
		return fail("Server-side error (update bus itinerary+reroute)")
	}

	// generate random uuid
	reqUuid := uuid.NewString()

	s.TrackingService.RegisterSession(tracker.TrackingSession{
		Id:       reqUuid,
		FromStop: *fromStop,
		ToStop:   *toStop,
		Bus:      bus,
	})

	return UserResponse{
		Ok:         true,
		Message:    "OK",
		TrackingID: reqUuid,
	}
}

func fail(msg string) UserResponse {
	return UserResponse{
		false,
		msg,
		"",
	}
}
