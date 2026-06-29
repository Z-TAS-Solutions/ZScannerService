#include "ControllerServer.h"
#include "ZLED.h"
#include <iostream>

ControllerServer::ControllerServer(std::shared_ptr<ZLED> indicator) : indicator_(indicator) {
    indicatorStatus_ = zscanproto::LEDStatus::VOID;
}

ControllerServer::~ControllerServer() {
}

grpc::Status ControllerServer::SetLEDStatus(grpc::ServerContext* context, const zscanproto::LEDStatusRequest* request, zscanproto::Status* response) {
    indicatorStatus_ = request->status();
    uint32_t r = 0, g = 0, b = 0;
    switch(request->status()) {
        case zscanproto::SUCCESS: g = 255; break;
        case zscanproto::PENDING: b = 255; break;
        case zscanproto::FAILED: r = 255; break;
        case zscanproto::VOID: default: break;
    }
    if (indicator_) indicator_->Set(r, g, b);

    response->set_success(true);
    response->set_message("LED Status Set successfully");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::GetLED(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::LEDState* response) {
    uint32_t r = 0, g = 0, b = 0;
    if (indicator_) indicator_->Get(r, g, b);

    response->set_red(r);
    response->set_green(g);
    response->set_blue(b);
    return grpc::Status::OK;
}

grpc::Status ControllerServer::ConfigureToF(grpc::ServerContext* context, const zscanproto::ToFConfig* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::EnableToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::DisableToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::ActivateToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::DeactivateToF(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::ActivateCamera(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::DeactivateCamera(grpc::ServerContext* context, const zscanproto::Empty* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::ConfigureCamera(grpc::ServerContext* context, const zscanproto::CameraConfig* request, zscanproto::Status* response) {
    response->set_success(false);
    response->set_message("Not implemented yet");
    return grpc::Status::OK;
}

grpc::Status ControllerServer::ToFEventStream(grpc::ServerContext* context, grpc::ServerReaderWriter<zscanproto::ToFEvent, zscanproto::ToFEvent>* stream) {
    return grpc::Status::OK;
}
