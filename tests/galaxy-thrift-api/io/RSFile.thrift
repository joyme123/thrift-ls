namespace java com.xiaomi.infra.galaxy.io.thrift
namespace php IO.RSFile
namespace py io.rsfile
namespace go io.rsfile

/***********************************************
 * RSFile(Record Sequence File)文件格式如下：
 *
 * | TCompactProtocol.encode(RSFileHeader) |
 * | Compressed(File Body)                 |
 *
 * 其中File Body格式为：
 *
 * | TCompactProtocol.encode(Record 1)   |
 * | TCompactProtocol.encode(Record 2)   |
 * | ...                                 |
 * | TCompactProtocol.encode(Record n)   |
 * | TCompactProtocol.encode(EOF Record) |
 *
 ***********************************************/

/**
 * 压缩算法类型
 */
enum Compression {
  NONE = 0,
  SNAPPY = 1,
}

/**
 * 数据完整性校验算法类型
 */
enum Checksum {
  NONE = 0,
  CRC32 = 1,
  // MD5
}

const string MAGIC = "RSF"

/**
 * 文件头
 */
struct RSFileHeader {
  /**
   * Magic常量，固定为"RSF"
   */
  1: optional string magic,
  /**
   * 版本号
   */
  2: optional i32 version,
  /**
   * 压缩算法类型
   * 可选，进行压缩的区域为文件中除文件头之外其余部分
   */
  3: optional Compression compression,
  /**
   * 数据完整性校验算法类型
   * 可选，对每条记录对data部分进行校验
   */
  4: optional Checksum checksum,
  /**
   * 记录数目，不含最后EOF记录，可选
   */
  5: optional i64 count = -1,
  /**
   * 元信息，内容由使用者自己定义
   */
  6: optional binary metadata,
}

/**
 * 记录
 */
struct Record {
  /**
   * 记录数据部分，仅用于非EOF记录
   */
  1: optional binary data,
  /**
   * 数据部分的校验值，可选
   */
  2: optional i32 checksum,
  /**
   * EOF记录，用于标记文件结束
   */
  3: optional bool eof = false,
}
