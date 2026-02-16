from concurrent import futures

import grpc
import grpc_gen
import routing.grpc_gen.optimization_pb2_grpc
from routing.grpc_gen.optimization_pb2_grpc import OptimizationServicer
import routing.grpc_gen.optimization_pb2 as optimization_pb2
from routing.optimize_functions import optimize_itinerary


# entrypoint for grpc optimization service
class OptimizeServer(OptimizationServicer):
    def RouteOptimize(self,request,context):
        optimize_route = optimize_itinerary()

        return optimization_pb2.OptimizedRoute(stops=[])

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    routing.grpc_gen.optimization_pb2_grpc.add_OptimizationServicer_to_server(OptimizeServer(),server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
