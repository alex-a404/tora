import json
from concurrent import futures
from time import sleep

import grpc
import numpy as np

import grpc_gen
import routing.grpc_gen.optimization_pb2_grpc
from routing.grpc_gen.optimization_pb2_grpc import OptimizationServicer
import routing.grpc_gen.optimization_pb2 as optimization_pb2
from routing.optimize_functions import optimize_itinerary
from google.protobuf.json_format import MessageToDict

ELEFTHERIA_STOP_ID = 581

# entrypoint for grpc optimization service
class OptimizeServer(OptimizationServicer):

    def __init__(self):
        self.dist_matrix = load_dist_mat()

    def RouteOptimize(self, request, context):
        try:
            req_dict = MessageToDict(request)

            constraints = []
            if 'constraints' in req_dict:
                for constraint in req_dict['constraints']:
                    constraints.append((constraint['From'], constraint['To']))

            # FIX 1: Pass the ENTIRE row for Eleftheria Square.
            # Now node_depths[581] works perfectly.
            node_depths = self.dist_matrix[ELEFTHERIA_STOP_ID]

            new_request = (req_dict['newRequest']['From'], req_dict['newRequest']['To'])

            optimized_route = optimize_itinerary(
                dist_matrix=np.array(self.dist_matrix),
                node_depths=node_depths,
                current_route=req_dict.get('itinerary', []),
                new_request=new_request,
                dependencies=constraints,
                direction=1,
            )
            print(f"Optimized: {optimized_route}")
            return optimization_pb2.OptimizedRoute(
                stops=optimized_route,
                success=True
            )

        except Exception as e:
            print(e)
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return optimization_pb2.OptimizedRoute(
                stops=[],
                success=False
            )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    routing.grpc_gen.optimization_pb2_grpc.add_OptimizationServicer_to_server(OptimizeServer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()

def load_dist_mat():
    with open("../data/distance_matrix.json", "r") as f:
        return json.load(f)

if __name__ == '__main__':
    serve()
