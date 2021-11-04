package chaincode

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-private-chaincode/ecc_go/chaincode/enclave"
	"github.com/hyperledger/fabric-private-chaincode/ecc_go/chaincode/ercc"
	"github.com/hyperledger/fabric-private-chaincode/ecc_go/chaincode/fakes"
	"github.com/hyperledger/fabric-private-chaincode/internal/crypto"
	"github.com/hyperledger/fabric-private-chaincode/internal/endorsement"
	"github.com/hyperledger/fabric-private-chaincode/internal/protos"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

//go:generate counterfeiter -o fakes/enclave.go -fake-name EnclaveStub . enclaveStub
//lint:ignore U1000 This is just used to generate fake
type enclaveStub interface {
	enclave.StubInterface
}

//go:generate counterfeiter -o fakes/utils.go -fake-name Extractors . extractors
//lint:ignore U1000 This is just used to generate fake
type extractors interface {
	Extractors
}

//go:generate counterfeiter -o fakes/validation.go -fake-name Validator . validator
//lint:ignore U1000 This is just used to generate fake
type validator interface {
	endorsement.Validation
}

//go:generate counterfeiter -o fakes/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
//lint:ignore U1000 This is just used to generate fake
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o fakes/ercc.go -fake-name ErccStub . erccStub
//lint:ignore U1000 This is just used to generate fake
type erccStub interface {
	ercc.Stub
}

func newECC(ec *enclave.EnclaveStub, val *fakes.Validator, ex *fakes.Extractors, ercc *fakes.ErccStub) *EnclaveChaincode {
	return &EnclaveChaincode{
		Enclave:   ec,
		Validator: val,
		Extractor: ex,
		Ercc:      ercc,
	}
}

func newFakes() (*fakes.EnclaveStub, *fakes.Validator, *fakes.Extractors, *fakes.ErccStub) {
	return &fakes.EnclaveStub{}, &fakes.Validator{}, &fakes.Extractors{}, &fakes.ErccStub{}
}

func newRealEc() *enclave.EnclaveStub {
	return enclave.NewEnclaveStub()
}

func TestInitECC(t *testing.T) {
	_, val, ex, ercc := newFakes()
	ecc := newECC(newRealEc(), val, ex, ercc)
	stub := &fakes.ChaincodeStub{}

	// test init
	r := ecc.Init(stub)
	assert.Equal(t, shim.Success(nil), r)

	// test invalid invocation
	stub.GetFunctionAndParametersReturns("whatever", nil)
	r = ecc.Invoke(stub)
	assert.Equal(t, shim.Error("invalid invocation"), r)
}

func TestEnclave(t *testing.T) {
	stub := &fakes.ChaincodeStub{}
	stub.GetFunctionAndParametersReturns("__initEnclave", nil)
	_, _, ex, _ := newFakes()
	ecc := newECC(newRealEc(), nil, ex, nil)

	attestParams := []byte("someAttestationParams")
	ccParams := &protos.CCParameters{
		ChaincodeId: "SomeChaincodeId",
	}
	hostParams := &protos.HostParameters{
		PeerMspId:    "",
		PeerEndpoint: "",
		Certificate:  nil,
	}

	ex.GetInitEnclaveMessageReturns(&protos.InitEnclaveMessage{AttestationParams: attestParams}, nil)
	ex.GetChaincodeParamsReturns(ccParams, nil)
	ex.GetHostParamsReturns(hostParams, nil)

	r := ecc.Invoke(stub)
	assert.EqualValues(t, shim.OK, r.Status)
	payload, err := base64.StdEncoding.DecodeString(string(r.Payload))
	assert.NoError(t, err)

	credentials := &protos.Credentials{}
	proto.Unmarshal(payload, credentials)

	stub.GetFunctionAndParametersReturns("__invoke", nil)
	ep := &crypto.EncryptionProviderImpl{
		CSP: crypto.GetDefaultCSP(),
		GetCcEncryptionKey: func() ([]byte, error) {
			attestedData := &protos.AttestedData{}
			err := proto.Unmarshal(credentials.SerializedAttestedData.GetValue(), attestedData)
			return []byte(base64.StdEncoding.EncodeToString(attestedData.GetChaincodeEk())), err
		}}

	ctx, _ := ep.NewEncryptionContext()
	requestBytes, _ := ctx.Conceal("invoke", []string{"a", "b", "10"})
	request, _ := base64.StdEncoding.DecodeString(requestBytes)
	chaincodeRequestMessage := &protos.ChaincodeRequestMessage{}
	err = proto.Unmarshal(request, chaincodeRequestMessage)
	ex.GetSerializedChaincodeRequestReturns(request, nil)

	//ec.ChaincodeInvokeReturns(expectedResp, nil)
	r = ecc.Invoke(stub)
	assert.EqualValues(t, shim.OK, r.Status)
	p, err := base64.StdEncoding.DecodeString(string(r.Payload))
	assert.NoError(t, err)
	//assert.EqualValues(t, expectedResp, p)
	//s, scr := ec.ChaincodeInvokeArgsForCall(1)
	//assert.Equal(t, stub, s)
	//assert.Equal(t, []byte("someChaincodeRequest"), scr)
	fmt.Println(p)

}

func expectError(t *testing.T, errorMsg string, r peer.Response) {
	assert.EqualValues(t, shim.ERROR, r.Status)
	assert.EqualValues(t, errorMsg, r.Message)
}
