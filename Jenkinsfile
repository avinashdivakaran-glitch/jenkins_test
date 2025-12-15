pipeline {
    agent any

    environment {
        APP_NAME = "tnn_backend"
        IMAGE_VERSION = "v1.0.0"

        SERVICES_BLUETOOTH = "service-bluetooth"
        SERVICES_SENSORS = "service-blesensors"
        SERVICES_MQTT = "service-mqtt"
        SERVICES_WIFI = "service-wifi"

        OCI_BUNDLE_DIR = "oci_bundles"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install Build Tools') {
            steps {
                sh """
                if ! command -v podman >/dev/null 2>&1; then
                    sudo apt update
                    sudo apt install -y podman golang
                fi
                # Ensure dpkg-deb is available (usually in build-essential or dpkg-dev)
                if ! command -v dpkg-deb >/dev/null 2>&1; then
                    sudo apt install -y dpkg-dev
                fi
                """
            }
        }

        stage('Load Pre-downloaded Base Podman Images') {
            steps {
                script {
                    sh """
                    cd /opt/jenkins_packages
                    sudo podman load -i debian-bookworm-slim.tar
                    sudo podman load -i golang-1.25.1.tar
                    """
                }
            }
        }

        stage('Build Podman Images') {
            steps {
                script {
                    sh """
                    sudo podman build --platform linux/arm64 -t ${SERVICES_BLUETOOTH}:${IMAGE_VERSION} ./service-bluetooth
                    sudo podman build --platform linux/arm64 -t ${SERVICES_SENSORS}:${IMAGE_VERSION} ./service-blesensors

                    sudo podman build --platform linux/arm64 -t ${SERVICES_WIFI}:${IMAGE_VERSION} ./service-wifi

                    sudo podman build --platform linux/arm64 -t ${SERVICES_MQTT}:${IMAGE_VERSION} ./service-mqtt
                    """
                }
            }
        }

        stage('Save Container Images') {
            steps {
                script {
                    sh """
                    sudo podman save -o service-bluetooth_${IMAGE_VERSION}.tar ${SERVICES_BLUETOOTH}:${IMAGE_VERSION}
                    sudo podman save -o service-blesensors_${IMAGE_VERSION}.tar ${SERVICES_SENSORS}:${IMAGE_VERSION}

                    sudo podman save -o service-wifi_${IMAGE_VERSION}.tar ${SERVICES_WIFI}:${IMAGE_VERSION}

                    sudo podman save -o service-mqtt_${IMAGE_VERSION}.tar ${SERVICES_MQTT}:${IMAGE_VERSION}
                    ls -lah
                    """
                }
            }
        }

        stage('Create OCI Bundle') {
            steps {
                script {
                    // Create OCI bundle directory if it doesn't exist
                    sh """
                    mkdir -p ${OCI_BUNDLE_DIR}
                    """

                    // Generate OCI bundles for each service
                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_BLUETOOTH}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_BLUETOOTH}
                    sudo podman push ${SERVICES_BLUETOOTH}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_BLUETOOTH}
                    // sudo tar -czvf oci_bundles/service-bluetooth.tar.gz -C oci_bundles service-bluetooth
                    """

                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_SENSORS}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_SENSORS}
                    sudo podman push ${SERVICES_SENSORS}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_SENSORS}
                    // sudo tar -czvf oci_bundles/service-blesensors.tar.gz -C oci_bundles service-blesensors
                    """

                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_WIFI}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_WIFI}
                    sudo podman push ${SERVICES_WIFI}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_WIFI}
                    // sudo tar -czvf oci_bundles/service-wifi.tar.gz -C oci_bundles service-wifi
                    """

                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_MQTT}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_MQTT}
                    sudo podman push ${SERVICES_MQTT}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_MQTT}
                    // sudo tar -czvf oci_bundles/service-mqtt.tar.gz -C oci_bundles service-mqtt
                    """
                }
            }
        }



        stage('Build Debian Package') {
            steps {
                script {
                    def pkgDir = "deb_temp"
                    def installPath = "/opt/${APP_NAME}/bundles"
                    
                    // 1. Clean and Create Directory Structure
                    sh """
                    rm -rf ${pkgDir}
                    mkdir -p ${pkgDir}/DEBIAN
                    mkdir -p ${pkgDir}${installPath}

                    mkdir -p ${pkgDir}/etc/containers/systemd
                    """

                    // 2. Copy OCI Artifacts (The images)
                    sh "sudo cp -r ${OCI_BUNDLE_DIR}/* ${pkgDir}${installPath}/"

                    // 3. Copy & Configure Control File
                    sh "cp packaging/control ${pkgDir}/DEBIAN/control"
                    sh "sed -i 's/VERSION_PLACEHOLDER/${IMAGE_VERSION}/g' ${pkgDir}/DEBIAN/control"
                    sh "sed -i 's/ARCH_PLACEHOLDER/${DEB_ARCH}/g' ${pkgDir}/DEBIAN/control"

                    // 4. Copy & Configure Quadlet Files
                    // Copy from repo to the temporary build folder
                    sh "cp packaging/quadlets/*.container ${pkgDir}/etc/containers/systemd/"
                    
                    // Replace the version placeholder in ALL container files at once
                    sh "sed -i 's/VERSION_PLACEHOLDER/${IMAGE_VERSION}/g' ${pkgDir}/etc/containers/systemd/*.container"

                    // 5. Copy Post-Install Script
                    sh "cp packaging/postinst.sh ${pkgDir}/DEBIAN/postinst"
                    sh "chmod 755 ${pkgDir}/DEBIAN/postinst"

                    // 6. Build the Package
                    sh "dpkg-deb --build ${pkgDir} ${DIST_DIR}/${APP_NAME}_${IMAGE_VERSION}_${DEB_ARCH}.deb"

                }
            }
        }



    }

    post {
        success {
            archiveArtifacts artifacts: '*.tar', fingerprint: true
            archiveArtifacts artifacts: '*.deb', fingerprint: true
            echo "Build successful! Download your .deb package from artifacts."
        }

        failure {
            // Send failure notification (optional)
            echo "Pipeline failed. Check the logs."
        }

        cleanup {
            // optional cleanup AFTER archiving
            cleanWs()
        }
    }
}


// sudo podman run -it   --name service-mqtt   --net=host   --privileged   -v /var/run/dbus:/var/run/dbus -v /usr/local/share/certs/broker:/certs:ro   service-mqtt:v1.0.0
// sudo podman run -it   --name service-bluetooth --net=host   --privileged   -v /var/run/dbus:/var/run/dbus -v /usr/local/share/certs/tnn-ble_app:/app/certs:ro   service-bluetooth:v1.0.0