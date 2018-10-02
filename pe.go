package main

import (
    "debug/pe"
    "encoding/binary"
)

type ClrHeader struct {
    IsDll                      bool
    HeaderSize                 uint32
    MajorRuntimeVersion        uint16
    MinorRuntimeVersion        uint16
    MetaDataDirectoryAddress   uint32
    MetaDataDirectorySize      uint32
    Flags                      uint32
    EntryPointToken            []byte
    ResourcesDirectoryAddress  uint32
    ResourcesDirectorySize     uint32
    StrongNameSignatureAddress uint32
    StrongNameSignatureSize    uint32
}

func readAssemblyMetadata(path string) (*ClrHeader, error) {
    f, err := pe.Open(path)

    if err != nil {
        return nil, err
    }
    defer f.Close()

    result := ClrHeader{}

    if (f.Characteristics & 0x2000) == 0x2000 {
        result.IsDll = true
    } else if (f.Characteristics & 0x2000) != 0x2000 {
        result.IsDll = false
    }

    txt := f.Section(".text")
    textBinary, _ := txt.Data()

    result.HeaderSize = binary.LittleEndian.Uint32(textBinary[8:12])
    result.MajorRuntimeVersion = binary.LittleEndian.Uint16(textBinary[12:14])
    result.MinorRuntimeVersion = binary.LittleEndian.Uint16(textBinary[14:16])
    result.MetaDataDirectoryAddress = binary.LittleEndian.Uint32(textBinary[16:20])
    result.MetaDataDirectorySize = binary.LittleEndian.Uint32(textBinary[20:24])
    result.Flags = binary.LittleEndian.Uint32(textBinary[24:28])
    result.EntryPointToken = textBinary[28:32]
    result.ResourcesDirectoryAddress = binary.LittleEndian.Uint32(textBinary[32:36])
    result.ResourcesDirectorySize = binary.LittleEndian.Uint32(textBinary[36:40])
    result.StrongNameSignatureAddress = binary.LittleEndian.Uint32(textBinary[40:44])
    result.StrongNameSignatureSize = binary.LittleEndian.Uint32(textBinary[44:48])

    return &result, nil
}
