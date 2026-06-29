cmake_minimum_required(VERSION 3.16)

project(ZScannerService LANGUAGES CXX)

set(CONAN_GEN_DIR "${CMAKE_CURRENT_SOURCE_DIR}/build/generators")

if(EXISTS "${CONAN_GEN_DIR}/conan_toolchain.cmake")
    include("${CONAN_GEN_DIR}/conan_toolchain.cmake")
    list(APPEND CMAKE_PREFIX_PATH "${CONAN_GEN_DIR}")
endif()

set(CMAKE_FIND_PACKAGE_PREFER_CONFIG ON)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

find_package(gRPC CONFIG REQUIRED)
find_package(protobuf CONFIG REQUIRED)

set(GEN_DIR "${CMAKE_BINARY_DIR}/generated")
file(MAKE_DIRECTORY ${GEN_DIR})

set(PROTO_SRC "${GEN_DIR}/ZScan.pb.cc")
set(PROTO_HDR "${GEN_DIR}/ZScan.pb.h")
set(GRPC_SRC  "${GEN_DIR}/ZScan.grpc.pb.cc")
set(GRPC_HDR  "${GEN_DIR}/ZScan.grpc.pb.h")
set(PROTO_FILES "${CMAKE_CURRENT_SOURCE_DIR}/src/proto/ZScan.proto")

add_custom_command(
    OUTPUT ${PROTO_SRC} ${PROTO_HDR} ${GRPC_SRC} ${GRPC_HDR}
    COMMAND $<TARGET_FILE:protobuf::protoc>
            --cpp_out=${GEN_DIR}
            --grpc_out=${GEN_DIR}
            --plugin=protoc-gen-grpc=$<TARGET_FILE:gRPC::grpc_cpp_plugin>
            -I ${CMAKE_CURRENT_SOURCE_DIR}/src/proto
            ${PROTO_FILES}
    DEPENDS ${PROTO_FILES}
    WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
    VERBATIM
)

add_custom_target(proto_gen ALL DEPENDS ${PROTO_SRC} ${PROTO_HDR} ${GRPC_SRC} ${GRPC_HDR})

if (CMAKE_SIZEOF_VOID_P EQUAL 8)
    set(ARCH x64)
else()
    set(ARCH x86)
endif()

file(GLOB EXTERNAL_LIBS "${CMAKE_CURRENT_SOURCE_DIR}/external/lib/Release/${ARCH}/*.lib")
file(GLOB EXTERNAL_ULIBS "${CMAKE_CURRENT_SOURCE_DIR}/external/lib/Universal/${ARCH}/*.lib")

set(APP_SOURCES
    src/main.cpp
    src/ControllerServer.cpp
    src/ZLED.cpp
)

add_executable(ZScannerService ${APP_SOURCES} ${PROTO_SRC} ${GRPC_SRC})
add_dependencies(ZScannerService proto_gen)

target_compile_definitions(ZScannerService PUBLIC 
    $<$<CONFIG:Debug>:ZSCAN_BUILD_DEBUG>
    $<$<CONFIG:Release>:ZSCAN_BUILD_RELEASE>
)
    

set(DLL_SOURCE_DIR "${CMAKE_CURRENT_SOURCE_DIR}/external/dll/$<CONFIG>/${ARCH}")

add_custom_command(
    TARGET ZScannerService 
    POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy_directory
        "${DLL_SOURCE_DIR}"
        "$<TARGET_FILE_DIR:ZScannerService>"
    COMMENT "Copying ZTAS Runtime DLLs for ${ARCH} ($<CONFIG>) to build directory..."
)


target_include_directories(ZScannerService PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${GEN_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/external/include
    ${CMAKE_CURRENT_SOURCE_DIR}/external/include/opencv2
)

target_link_libraries(ZScannerService PRIVATE
    gRPC::grpc++
    protobuf::libprotobuf
    gpiod
    ${EXTERNAL_ULIBS}
    ${EXTERNAL_LIBS}
)
