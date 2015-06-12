package chameleon

import (
	"time"
)

type Packet struct {
	RelTime time.Duration
	Length int64
}

type Conn struct {
	Cipher string
	Packets []Packet
}

type Signature struct {
	Conns []Conn
}

func (s *Conn) Dimension() time.Duration, time.Duration, int64 {
	mint, maxt := time.Duration(0)
	maxl := 0
	
	for i, pkt := range(s.Packets) {
		if i == 0 || pkt.RelTime < min {
			mint = pkt.RelTime
		}
		
		if i == 0 || pkt.RelTime > max {
			maxt = pkt.RelTime
		}
		
		if pkt.Length > maxl {
			maxl = pkt.Length
		}
	}
	
	return mint, maxt, maxl
}

func (s *Conn) Graph(rest, resl int64) {
	mint, maxt, maxl := s.Dimension()
	
	result := make([][]int, (maxt - mint) / rest)
	
	for t = 0; t < len(result); t++ {
		result[t] = make([]int, maxl / resl)
	}
	
	for _, pkt := range s.Packets {
		result[(pkt.RelTime - mint) / rest][pkt.Length / resl] += 1
	}
	
	return result
}

type sigStreamFactory struct {
	Signature
}

func (sf *sigStreamFactory) New() tcpassembly.Stream {
	conn := make(Conn)
	sf.Signature.Conns = append(sf.Signature.Conns, conn)
	
	r := tcpreader.NewReaderStream()
	
	go func() {
		header := make([]byte, 5)
		
		for {
			_, err := io.ReadFull(r, header)
			if err != nil {
				panic(err)
			}
			
			lng := r[3] << 8 | r[4]
			
			log.Info("Type: %d\n", r[0])
			log.Info("Length: %d\n", lng)
			
			io.CopyN(ioutil.Discard, r, lng)
		}
	}()
	
	return &r
}

func SignatureFromCapture(h *pcap.Handle) Signature, error {
	ps := gopacket.NewPacketSource(h, h.LinkType())
	pkts := ps.Packets()
	
	sig := &Signature{
		Conns: make([]Conn),
	}
	
	sf := &sigStreamFactory{
		Signature: sig,
	}
	sp := tcpassembly.NewStreamPool(sf)
	as := tcpassembly.NewAssembler(sp)
	
	for pkt := range pkts {
		if pkt == nil {
			break
		}
		
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil {
			continue
		}
		
		if packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}
		
		tcp := pkt.TransportLayer().(*layers.TCP)
		as.AssembleWithTimestamp(pkt.NetworkLayer().NetworkFlow(), tcp,
			pkt.Metadata().Timestamp)
	}
	
	assembler.FlushAll()
	
	return sig
}