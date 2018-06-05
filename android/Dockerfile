FROM openjdk:8 AS android
ENV SDK_URL="https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip" \
    ANDROID_HOME="/usr/local/android-sdk"
RUN mkdir "${ANDROID_HOME}" .android
WORKDIR "${ANDROID_HOME}"
RUN curl -o sdk.zip "${SDK_URL}" && unzip sdk.zip && rm sdk.zip
ENV ANDROID_VERSION=26 \
    ANDROID_BUILD_TOOLS_VERSION=27.0.3 \
    PATH=${PATH}:${ANDROID_HOME}/tools:${ANDROID_HOME}/tools/bin:${ANDROID_HOME}/platform-tools
RUN yes | ${ANDROID_HOME}/tools/bin/sdkmanager --licenses
RUN echo suppress warning > /root/.android/repositories.cfg
RUN sdkmanager "build-tools;${ANDROID_BUILD_TOOLS_VERSION}" \
               "platforms;android-${ANDROID_VERSION}" \
               "platform-tools"

FROM android as androidbuilder
COPY android /application/
VOLUME /outputs
RUN echo sdk.dir=${ANDROID_HOME} > /application/local.properties
WORKDIR /application
RUN ./gradlew clean && ./gradlew assembleRelease
CMD cp -R /application/app/build/outputs/* /outputs && chown -R $USER:$(id -g) /outputs
