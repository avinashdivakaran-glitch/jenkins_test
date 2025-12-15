pipeline {
    agent any

    environment {
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

        stage('Install Podman') {
            steps {
                sh """
                if ! command -v podman >/dev/null 2>&1; then
                    sudo apt update
                    sudo apt install -y podman golang
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
                    sudo tar -czvf oci_bundles/service-bluetooth.tar.gz -C oci_bundles service-bluetooth
                    """

                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_SENSORS}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_SENSORS}
                    sudo podman push ${SERVICES_SENSORS}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_SENSORS}
                    sudo tar -czvf oci_bundles/service-blesensors.tar.gz -C oci_bundles service-blesensors
                    """

                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_WIFI}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_WIFI}
                    sudo podman push ${SERVICES_WIFI}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_WIFI}
                    sudo tar -czvf oci_bundles/service-wifi.tar.gz -C oci_bundles service-wifi
                    """

                    sh """
                    sudo rm -rf ${OCI_BUNDLE_DIR}/${SERVICES_MQTT}
                    mkdir -p ${OCI_BUNDLE_DIR}/${SERVICES_MQTT}
                    sudo podman push ${SERVICES_MQTT}:${IMAGE_VERSION} oci:${OCI_BUNDLE_DIR}/${SERVICES_MQTT}
                    sudo tar -czvf oci_bundles/service-mqtt.tar.gz -C oci_bundles service-mqtt
                    """
                }
            }
        }

    }

    post {
        success {
            archiveArtifacts artifacts: '*.tar', fingerprint: true
            archiveArtifacts artifacts: 'oci_bundles/*.tar.gz', fingerprint: true
            echo "Pipeline successfully completed!"
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