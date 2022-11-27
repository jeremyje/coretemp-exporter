# Copyright 2022 Jeremy Edwards
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM --platform=$BUILDPLATFORM gcr.io/distroless/static:nonroot

ARG BINARY_PATH

WORKDIR /app/

COPY --chown=nonroot $BINARY_PATH /app/example
COPY --chown=nonroot LICENSE /app/LICENSE

EXPOSE 80
ENTRYPOINT ["/app/example"]

# Standardized Docker Image Labels
# https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL \
  org.opencontainers.image.created="${BUILD_TIME}" \
  org.opencontainers.image.authors="Jeremyje <jeremyje@googlegroups.com>" \
  org.opencontainers.image.url="https://github.com/jeremyje/gomain" \
  org.opencontainers.image.documentation="https://github.com/jeremyje/gomain" \
  org.opencontainers.image.source="https://github.com/jeremyje/gomain" \
  org.opencontainers.image.version="${BUILD_VERSION}" \
  org.opencontainers.image.revision="1" \
  org.opencontainers.image.vendor="Jeremy Edwards" \
  org.opencontainers.image.licenses="Apache-2.0" \
  org.opencontainers.image.ref.name="" \
  org.opencontainers.image.title="${BINARY_NAME}" \
  org.opencontainers.image.description="Gomain Example App" \
  org.label-schema.schema-version="1.0" \
  org.label-schema.build-date=$BUILD_DATE \
  org.label-schema.url="https://github.com/jeremyje/gomain" \
  org.label-schema.vcs-url="https://github.com/jeremyje/gomain" \
  org.label-schema.version=$BUILD_VERSION \
  org.label-schema.vcs-ref=$VCS_REF \
  org.label-schema.vendor="Jeremy Edwards" \
  org.label-schema.name="${BINARY_NAME}" \
  org.label-schema.description="Gomain Example App" \
  org.label-schema.usage="https://github.com/jeremyje/gomain"
