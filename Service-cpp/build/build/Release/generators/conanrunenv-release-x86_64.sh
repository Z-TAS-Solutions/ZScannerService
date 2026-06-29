script_folder="/mnt/d/Workspace/Repos/Z-TAS-Solutions/ZScannerService/Service-cpp/build/build/Release/generators"
echo "echo Restoring environment" > "$script_folder/deactivate_conanrunenv-release-x86_64.sh"
for v in GRPC_DEFAULT_SSL_ROOTS_FILE_PATH OPENSSL_MODULES
do
   is_defined="true"
   value=$(printenv $v) || is_defined="" || true
   if [ -n "$value" ] || [ -n "$is_defined" ]
   then
       echo export "$v='$value'" >> "$script_folder/deactivate_conanrunenv-release-x86_64.sh"
   else
       echo unset $v >> "$script_folder/deactivate_conanrunenv-release-x86_64.sh"
   fi
done

export GRPC_DEFAULT_SSL_ROOTS_FILE_PATH="/home/zischl/.conan2/p/b/grpc85988c062d451/p/res/grpc/roots.pem"
export OPENSSL_MODULES="/home/zischl/.conan2/p/b/opens6233d39fdded6/p/lib/ossl-modules"