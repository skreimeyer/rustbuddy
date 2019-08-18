
#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

oses=("windows" "darwin" "linux")
arches=("amd64" "386")

for os in "${oses[@]}"
do
    for arch in "${arches[@]}"
    do
        output_name=$package_name'-'$os'-'$arch
        if [ $os = "windows" ]; then
            output_name+='.exe'
        fi  

        env GOOS=$os GOARCH=$arch go build -o $output_name $package
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
        fi
    done
done