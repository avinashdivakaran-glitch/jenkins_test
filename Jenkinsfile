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
                    tar -czvf ${OCI_BUNDLE_DIR}/service-bluetooth.tar.gz -C ${OCI_BUNDLE_DIR}/${SERVICES_BLUETOOTH}

                    """
                }
            }
        }

    }

    post {
        success {
            archiveArtifacts artifacts: '*.tar', fingerprint: true
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
