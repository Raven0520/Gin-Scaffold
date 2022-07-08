package util

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

// TraceLogic 追踪结构体
type TraceBody struct {
	TraceID     string // 链路ID
	SpanID      string // Span ID
	Caller      string //
	SrcMethod   string //
	HintCode    int64  //
	HintContent string //
}

// TContext Trace 上下文
type TContext struct {
	TraceBody
	CSpanID string
}

// NewTrace 实例化 Trace
func NewTrace() *TContext {
	trace := &TContext{}
	trace.TraceID = GetTraceID()
	trace.SpanID = NewSpanID()
	return trace
}

// NewSpanID 实例化 SpanID
func NewSpanID() string {
	timestamp := uint32(time.Now().Unix())
	ipToLong := binary.BigEndian.Uint32(LocalIP.To4())
	b := bytes.Buffer{}
	b.WriteString(fmt.Sprintf("%08x", ipToLong^timestamp))
	b.WriteString(fmt.Sprintf("%08x", rand.Int31()))
	return b.String()
}

// 计算 TraceID
func calcTraceID(ip string) (traceID string) {
	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()

	b := bytes.Buffer{}
	netIP := net.ParseIP(ip)
	if netIP == nil {
		b.WriteString("00000000")
	} else {
		b.WriteString(hex.EncodeToString(netIP.To4()))
	}
	b.WriteString(fmt.Sprintf("%08x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0") // 末两位标记来源,b0为go

	return b.String()
}

// GetTraceId 根据IP地址 计算链路ID
func GetTraceID() (traceID string) {
	return calcTraceID(LocalIP.String())
}
