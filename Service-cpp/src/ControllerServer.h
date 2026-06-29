#pragma once

#include <grpcpp/grpcpp.h>
#include "ZScan.grpc.pb.h"
#include <mutex>
#include <memory>

class ZLED;

class ControllerServer final : public zscanproto::ZPiController::Service {
public:
    ControllerServer(std::shared_ptr<ZLED> indicator);
    ~ControllerServer();

    grpc::Status SetLEDStatus(grpc::ServerContext* context, const zscanproto::LEDStatusRequest* request, zscanproto::Status* response) override;
    grpc::Status GetLED(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::LEDState* response) override;

    grpc::Status ConfigureToF(grpc::ServerContext* context, const zscanproto::ToFConfig* request, zscanproto::Status* response) override;
    grpc::Status EnableToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) override;
    grpc::Status DisableToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) override;
    grpc::Status ActivateToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) override;
    grpc::Status DeactivateToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) override;

    grpc::Status ActivateCamera(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) override;
    grpc::Status DeactivateCamera(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) override;
    grpc::Status ConfigureCamera(grpc::ServerContext* context, const zscanproto::CameraConfig* request, zscanproto::Status* response) override;

    grpc::Status ToFEventStream(grpc::ServerContext* context, grpc::ServerReaderWriter<zscanproto::ToFEvent, zscanproto::ToFEvent>* stream) override;

private:
    zscanproto::LEDStatus indicatorStatus_;
    std::shared_ptr<ZLED> indicator_;
};
