// Copyright 2019 by mauro@ezplanet.org (Mauro Mozzarelli)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
// this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation
// and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package base

const StrERROR = "ERROR"
const TcpPORT = "2224"

// Socketmap standard errors and responses
const SmapOK = "OK "
const SmapNOTFOUND = "NOTFOUND "
const SmapTEMP = "TEMP "
const SmapPERM = "PERM "
const StrOK = "OK"

const StrREJECT = "REJECT"
const StrKEEPALIVE = "KEEPALIVE"

// Generic errors
const StrTemporaryERROR = "temporary error, please try again later"
const StrInternalERROR = "internal configuration error, please try again later"
const StrTerminated = "program terminated"

// Netstring errors
const ErrNsInvalid = "invalid netstring format"
const ErrNsLenNotNumeric = "netstring length is not numeric"
const ErrNsLenMismatch = "netstring length mismatch"
const ErrNsEmpty = "empty netstring"

// Request Parser Errors
const ErrParseInvalid = "invalid request format"
