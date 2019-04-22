# coding: utf-8

import grpc
import cqssc_pb2
import cqssc_pb2_grpc

_ADDR = '0.0.0.0:3721'


def run_validate(stub):
    req = cqssc_pb2.ValidateReq(
        tag="x5.eq.batch", 
        code=["01", "12", "3", "4", "1235"]
    )
    try:
        resp = stub.Validate(req, timeout=3)
    except grpc.RpcError as e:
        print e
        return

    print "validate client received: ", resp.num

def run_compute(stub):
    req = cqssc_pb2.ComputeReq(
        tag="x5.eq.batch",
        result=["0", "2", "3", "4", "5"],
        code=["01", "12", "3", "4", "1235"],
    )
    try:
        resp = stub.Compute(req, timeout=3)

    except grpc.RpcError as e:
        print e
        return

    print "compute client received: ", resp.num




if __name__ == '__main__':
    with grpc.insecure_channel(_ADDR) as channel:
        stub = cqssc_pb2_grpc.FormulaStub(channel)
        run_validate(stub)
        run_compute(stub)