#!/bin/bash
#
# Copyright 2020 Aletheia Ware LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
set -x

(cd $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/ui/data/ && ./gen.sh)
go fmt $GOPATH/src/github.com/AletheiaWareLLC/{labclientgo,labclientgo/...}
go vet $GOPATH/src/github.com/AletheiaWareLLC/{labclientgo,labclientgo/...}
go test $GOPATH/src/github.com/AletheiaWareLLC/{labclientgo,labclientgo/...}
ANDROID_NDK_HOME=${ANDROID_HOME}/ndk-bundle/
(cd $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/cmd && fyne package -os android -appID com.aletheiaware.lab -icon $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/lab.svg -name LabAndroid)
#03-25 21:17:21.397 21848 21848 I letheiaware.la: Late-enabling -Xcheck:jni
#03-25 21:17:21.418 21848 21848 E letheiaware.la: Unknown bits set in runtime_flags: 0x8000
#03-25 21:17:21.448 21848 21848 W letheiaware.la: Can't mmap dex file /data/app/com.aletheiaware.lab-naC86cdO6jg6E9PTJnbWag==/base.apk!classes.dex directly; please zipalign to 4 bytes. Falling back to extracting file.
(cd $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/cmd && ${ANDROID_HOME}/build-tools/28.0.3/zipalign -f 4 LabAndroid.apk LabAndroid-aligned.apk)
(cd $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/cmd && adb install -r LabAndroid-aligned.apk)
#(cd $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/cmd && adb logcat com.aletheiaware.lab:V org.golang.app:V *:S | tee log)
(cd $GOPATH/src/github.com/AletheiaWareLLC/labclientgo/cmd && adb logcat -c && adb logcat | tee android.log)
