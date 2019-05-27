// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

const (
	defaultRpcAddress    = "tcp://localhost:46657"
	defaultAbis          = "[{\"constant\":false,\"inputs\":[],\"name\":\"add\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"get\",\"outputs\":[{\"name\":\"\",\"type\":\"int32\"}],\"payable\":false,\"type\":\"function\"}]"
	defaultBytecode      = "6060604052341561000f57600080fd5b5b6101058061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634f2be91f1460475780636d4ce63c146059575b600080fd5b3415605157600080fd5b60576085565b005b3415606357600080fd5b606960c2565b604051808260030b60030b815260200191505060405180910390f35b60008081819054906101000a900460030b8092919060010191906101000a81548163ffffffff021916908360030b63ffffffff160217905550505b565b60008060009054906101000a900460030b90505b905600a165627a7a72305820259a0a3f2a8a112df2232529a36c75cc314d05060713c663a0786913fee723160029"
	defaultContractAddr  = "0x0f3200148775219ead5ba8d2f2bf9f0ae1f76ea9"
	defaultToAccountAddr = "0x7752b42608a0f1943c19fc5802cb027e60b4c911"
	defaultPrivKey       = "7d73c3dafd3c0215b8526b26f8dbdb93242fc7dcfbdfa1000d93436d577c3b94"
	defaultChainID       = "annchain-genesis"
)
