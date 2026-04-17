#include <iostream>
#include <memory>
#include <string>
#include <grpcpp/grpcpp.h>
#include "ControllerServer.h"
#include "ZLED.h"

void RunServer() {
    std::string server_address("0.0.0.0:50051");
    auto indicator = std::make_shared<ZLED>(17, 27, 22);
    ControllerServer service(indicator);

    grpc::ServerBuilder builder;
    builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
    builder.RegisterService(&service);

    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    std::cout << "Server listening on " << server_address << std::endl;

    server->Wait();
}

int main(int argc, char** argv) {
    RunServer();
    return 0;
}
