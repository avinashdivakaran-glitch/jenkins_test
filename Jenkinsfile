pipeline {
    agent any

    environment {
        IMAGE_VERSION = "v1.0.0"

        SERVICES_BLUETOOTH = "service-bluetooth"
        SERVICES_SENSORS = "service-blesensors"
        SERVICES_MQTT = "service-mqtt"

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
                    cd /opt/jankins_packages
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
                    sudo podman build --platform linux/arm64 -t ${SERVICES_MQTT}:${IMAGE_VERSION} ./service-mqtt
                    """
                }
            }
        }

        // stage('Create OCI Bundle') {
        //     steps {
        //         script {
        //             // Create OCI bundle directory if it doesn't exist
        //             sh """
        //             mkdir -p ${OCI_BUNDLE_DIR}
        //             """

        //             // Generate OCI bundles for each service
        //             sh """
        //             sudo podman generate systemd --name ${SERVICES_BLUETOOTH}:${IMAGE_VERSION} --output ${OCI_BUNDLE_DIR}/service-bluetooth-bundle
        //             sudo podman generate systemd --name ${SERVICES_SENSORS}:${IMAGE_VERSION} --output ${OCI_BUNDLE_DIR}/service-sensors-bundle
        //             sudo podman generate systemd --name ${SERVICES_MQTT}:${IMAGE_VERSION} --output ${OCI_BUNDLE_DIR}/service-mqtt-bundle
        //             """
        //         }
        //     }
        // }

    }

    post {
        always {
            // Clean up any temporary files or artifacts
            cleanWs()
        }

        success {
            // Send success notification (optional)
            echo "Pipeline successfully completed!"
        }

        failure {
            // Send failure notification (optional)
            echo "Pipeline failed. Check the logs."
        }
    }
}
