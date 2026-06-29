script_folder="/mnt/d/Workspace/Repos/Z-TAS-Solutions/ZScannerService/Service-cpp/build/Release/generators"
echo "echo Restoring environment" > "$script_folder/deactivate_conanrunenv-release-armv8.sh"
for v in GRPC_DEFAULT_SSL_ROOTS_FILE_PATH OPENSSL_MODULES
do
   is_defined="true"
   value=$(printenv $v) || is_defined="" || true
   if [ -n "$value" ] || [ -n "$is_defined" ]
   then
       echo export "$v='$value'" >> "$script_folder/deactivate_conanrunenv-release-armv8.sh"
   else
       echo unset $v >> "$script_folder/deactivate_conanrunenv-release-armv8.sh"
   fi
done

export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH="/home/zischl/.conan2/p/b/grpc4c6b4e87a7adf/p/res/grpc/roots.pem"
export OPENSSL_MODULES="/home/zischl/.conan2/p/b/opens642a7e72256be/p/lib/ossl-modules"