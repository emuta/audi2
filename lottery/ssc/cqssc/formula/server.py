# coding: utf-8

"""

Document:
https://grpc.io/grpc/python/grpc.html#service-side-context

Example:
https://github.com/grpc/grpc/tree/master/examples/python


handler error:
context.set_code(grpc.StatusCode.UNIMPLEMENTED)
context.set_details('Method not implemented!')


"""

import time
from concurrent import futures
import grpc
import cqssc_pb2
import cqssc_pb2_grpc
import cqssc_validate
import cqssc_compute

_ADDR = '0.0.0.0:3721'
_ONE_DAY_IN_SECONDS = 60 * 60 * 24

class FormulaServicer(cqssc_pb2_grpc.FormulaServicer):

    def Validate(self, request, context):
        print request.tag, request.code

        num = cqssc_validate.Validate(request.tag, request.code)
        result = dict(num=num)

        return cqssc_pb2.ValidateResp(**result)

    def Compute(self, request, context):
        print request.tag, request.result, request.code

        num = cqssc_compute.Compute(request.tag, request.result, request.code)
        result = dict(num=num)
        return cqssc_pb2.ComputeResp(**result)

def main():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    cqssc_pb2_grpc.add_FormulaServicer_to_server(FormulaServicer(), server)

    server.add_insecure_port(_ADDR)
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS) # one day in seconds
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    main()