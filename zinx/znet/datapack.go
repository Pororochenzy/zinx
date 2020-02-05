package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"pororo.com/zinx/utils"
	"pororo.com/zinx/ziface"

)

//封包拆包类实例，暂时不需要成员
type DataPack struct {}

//封包拆包实例初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头长度方法
func(dp *DataPack) GetHeadLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}
//封包方法(压缩数据)
func(dp *DataPack) Pack(msg ziface.IMessage)([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{}) //相当于以[]byte为初始字节切片， 创建缓冲，缓冲流就是那些读完就没，那种

 //问题： 字节切片 里面不是已经是二进制了吗  例如[]byte{'a','b','c'} -> [97 98 99 ]， 或者说什么是二进制编码 ， 【也不知道它以什么码表 去转换成二进制】 实验打印看 是按照ASCII表
 //--》，ASCII将字母、数字和其它符号编号，并用7比特的二进制来表示这个整数，所以说看【97 98 99] 相当于二进制了

 /*
 func Write(w io.Writer, order ByteOrder, data interface{}) error
将data的binary编码格式写入w，data必须是定长值、定长值的切片、定长值的指针。order指定写入数据的字节序，写入结构体时，名字中有'_'的字段会置为0。
  */
	//写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {  //如果data参数是个'你',那么以二进制编码写进去的话 打印是 [96 79 0 0]
		return nil, err
	}
	fmt.Println("二进制编码格式 写入的是什么东西 ：",dataBuff.Bytes()) //打印一下 看看  二进制编码格式 写入的是什么东西// 打印：msg.GetDataLen()是5  ,-[5 0 0 0]
	//msg.GetDataLen()是7的话 打印-》  [7 0 0 0]

	//写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil ,err
	}
	//写入一个 会占4个字节  --》因为int32 所以也是占4个字节       //DataLen,MsgId,Data 的顺序
	fmt.Println("封包后，包里的数据 ：",dataBuff.Bytes()) // 7 ，1  ，[]byte{'w', 'o', 'r', 'l', 'd', '!', '!'}, -->打印[7 0 0 0 1 0 0 0 119 111 114 108 100 33 33]

	return dataBuff.Bytes(), nil
}
//拆包方法(解压数据)
//创建一个从输入二进制数据的ioReader
func(dp *DataPack) Unpack(binaryData []byte)(ziface.IMessage, error) {
	fmt.Println("拆包前二进制编码格式 是什么东西 ：",binaryData) //打印一下 看看

	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	//判断dataLen的长度是否超出我们允许的最大包长度
	if (utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen >utils.GlobalObject.MaxPacketSize) {
		return nil, errors.New("Too large msg data recieved")
	}

	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
/*需要注意的是整理的`Unpack`方法，因为我们从上图可以知道，我们进行拆包的时候是分两次过程的，
第二次是依赖第一次的dataLen结果，所以`Unpack`只能解压出包头head的内容，得到msgId 和  dataLen。之后调用者再根据dataLen继续从io流中读取body中的数据。
 */