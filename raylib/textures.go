package raylib

/*
#include "raylib.h"
#include <stdlib.h>
*/
import "C"

import (
	"image"
	"unsafe"
)

// TextureFormat - Texture format
type TextureFormat int32

// Texture formats
// NOTE: Support depends on OpenGL version and platform
const (
	// 8 bit per pixel (no alpha)
	UncompressedGrayscale TextureFormat = C.UNCOMPRESSED_GRAYSCALE
	// 16 bpp (2 channels)
	UncompressedGrayAlpha TextureFormat = C.UNCOMPRESSED_GRAY_ALPHA
	// 16 bpp
	UncompressedR5g6b5 TextureFormat = C.UNCOMPRESSED_R5G6B5
	// 24 bpp
	UncompressedR8g8b8 TextureFormat = C.UNCOMPRESSED_R8G8B8
	// 16 bpp (1 bit alpha)
	UncompressedR5g5b5a1 TextureFormat = C.UNCOMPRESSED_R5G5B5A1
	// 16 bpp (4 bit alpha)
	UncompressedR4g4b4a4 TextureFormat = C.UNCOMPRESSED_R4G4B4A4
	// 32 bpp
	UncompressedR8g8b8a8 TextureFormat = C.UNCOMPRESSED_R8G8B8A8
	// 4 bpp (no alpha)
	CompressedDxt1Rgb TextureFormat = C.COMPRESSED_DXT1_RGB
	// 4 bpp (1 bit alpha)
	CompressedDxt1Rgba TextureFormat = C.COMPRESSED_DXT1_RGBA
	// 8 bpp
	CompressedDxt3Rgba TextureFormat = C.COMPRESSED_DXT3_RGBA
	// 8 bpp
	CompressedDxt5Rgba TextureFormat = C.COMPRESSED_DXT5_RGBA
	// 4 bpp
	CompressedEtc1Rgb TextureFormat = C.COMPRESSED_ETC1_RGB
	// 4 bpp
	CompressedEtc2Rgb TextureFormat = C.COMPRESSED_ETC2_RGB
	// 8 bpp
	CompressedEtc2EacRgba TextureFormat = C.COMPRESSED_ETC2_EAC_RGBA
	// 4 bpp
	CompressedPvrtRgb TextureFormat = C.COMPRESSED_PVRT_RGB
	// 4 bpp
	CompressedPvrtRgba TextureFormat = C.COMPRESSED_PVRT_RGBA
	// 8 bpp
	CompressedAstc4x4Rgba TextureFormat = C.COMPRESSED_ASTC_4x4_RGBA
	// 2 bpp
	CompressedAstc8x8Rgba TextureFormat = C.COMPRESSED_ASTC_8x8_RGBA
)

// TextureFilterMode - Texture filter mode
type TextureFilterMode int32

// Texture parameters: filter mode
// NOTE 1: Filtering considers mipmaps if available in the texture
// NOTE 2: Filter is accordingly set for minification and magnification
const (
	// No filter, just pixel aproximation
	FilterPoint TextureFilterMode = C.FILTER_POINT
	// Linear filtering
	FilterBilinear TextureFilterMode = C.FILTER_BILINEAR
	// Trilinear filtering (linear with mipmaps)
	FilterTrilinear TextureFilterMode = C.FILTER_TRILINEAR
	// Anisotropic filtering 4x
	FilterAnisotropic4x TextureFilterMode = C.FILTER_ANISOTROPIC_4X
	// Anisotropic filtering 8x
	FilterAnisotropic8x TextureFilterMode = C.FILTER_ANISOTROPIC_8X
	// Anisotropic filtering 16x
	FilterAnisotropic16x TextureFilterMode = C.FILTER_ANISOTROPIC_16X
)

// TextureWrapMode - Texture wrap mode
type TextureWrapMode int32

// Texture parameters: wrap mode
const (
	WrapRepeat TextureWrapMode = C.WRAP_REPEAT
	WrapClamp  TextureWrapMode = C.WRAP_CLAMP
	WrapMirror TextureWrapMode = C.WRAP_MIRROR
)

// Image type, bpp always RGBA (32bit)
// NOTE: Data stored in CPU memory (RAM)
type Image struct {
	// Image raw data
	Data unsafe.Pointer
	// Image base width
	Width int32
	// Image base height
	Height int32
	// Mipmap levels, 1 by default
	Mipmaps int32
	// Data format (TextureFormat)
	Format TextureFormat
}

func (i *Image) cptr() *C.Image {
	return (*C.Image)(unsafe.Pointer(i))
}

// ToImage converts a Image to Go image.Image
func (i *Image) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, int(i.Width), int(i.Height)))

	// Get pixel data from image (RGBA 32bit)
	pixels := GetImageData(i)

	img.Pix = (*[1 << 30]uint8)(pixels)[:]

	return img
}

// NewImage - Returns new Image
func NewImage(data unsafe.Pointer, width, height, mipmaps int32, format TextureFormat) *Image {
	return &Image{data, width, height, mipmaps, format}
}

// NewImageFromPointer - Returns new Image from pointer
func NewImageFromPointer(ptr unsafe.Pointer) *Image {
	return (*Image)(ptr)
}

// NewImageFromImage - Returns new Image from Go image.Image
func NewImageFromImage(img image.Image) *Image {
	size := img.Bounds().Size()
	pixels := make([]Color, size.X*size.Y)

	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			color := img.At(x, y)
			r, g, b, a := color.RGBA()
			pixels[x+y*size.Y] = NewColor(uint8(r), uint8(g), uint8(b), uint8(a))
		}
	}

	return LoadImageEx(pixels, int32(size.X), int32(size.Y))
}

// Texture2D type, bpp always RGBA (32bit)
// NOTE: Data stored in GPU memory
type Texture2D struct {
	// OpenGL texture id
	ID uint32
	// Texture base width
	Width int32
	// Texture base height
	Height int32
	// Mipmap levels, 1 by default
	Mipmaps int32
	// Data format (TextureFormat)
	Format TextureFormat
}

func (t *Texture2D) cptr() *C.Texture2D {
	return (*C.Texture2D)(unsafe.Pointer(t))
}

// NewTexture2D - Returns new Texture2D
func NewTexture2D(id uint32, width, height, mipmaps int32, format TextureFormat) Texture2D {
	return Texture2D{id, width, height, mipmaps, format}
}

// NewTexture2DFromPointer - Returns new Texture2D from pointer
func NewTexture2DFromPointer(ptr unsafe.Pointer) Texture2D {
	return *(*Texture2D)(ptr)
}

// RenderTexture2D type, for texture rendering
type RenderTexture2D struct {
	// Render texture (fbo) id
	ID uint32
	// Color buffer attachment texture
	Texture Texture2D
	// Depth buffer attachment texture
	Depth Texture2D
}

func (r *RenderTexture2D) cptr() *C.RenderTexture2D {
	return (*C.RenderTexture2D)(unsafe.Pointer(r))
}

// NewRenderTexture2D - Returns new RenderTexture2D
func NewRenderTexture2D(id uint32, texture, depth Texture2D) RenderTexture2D {
	return RenderTexture2D{id, texture, depth}
}

// NewRenderTexture2DFromPointer - Returns new RenderTexture2D from pointer
func NewRenderTexture2DFromPointer(ptr unsafe.Pointer) RenderTexture2D {
	return *(*RenderTexture2D)(ptr)
}

// LoadImage - Load an image into CPU memory (RAM)
func LoadImage(fileName string) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadImage(cfileName)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageEx - Load image data from Color array data (RGBA - 32bit)
func LoadImageEx(pixels []Color, width, height int32) *Image {
	cpixels := pixels[0].cptr()
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.LoadImageEx(cpixels, cwidth, cheight)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImagePro - Load image from raw data with parameters
func LoadImagePro(data []byte, width, height int32, format TextureFormat) *Image {
	cdata := unsafe.Pointer(&data[0])
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	ret := C.LoadImagePro(cdata, cwidth, cheight, cformat)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadImageRaw - Load image data from RAW file
func LoadImageRaw(fileName string, width, height int32, format TextureFormat, headerSize int32) *Image {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cformat := (C.int)(format)
	cheaderSize := (C.int)(headerSize)
	ret := C.LoadImageRaw(cfileName, cwidth, cheight, cformat, cheaderSize)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadTexture - Load an image as texture into GPU memory
func LoadTexture(fileName string) Texture2D {
	cfileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cfileName))
	ret := C.LoadTexture(cfileName)
	v := NewTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadTextureFromImage - Load a texture from image data
func LoadTextureFromImage(image *Image) Texture2D {
	cimage := image.cptr()
	ret := C.LoadTextureFromImage(*cimage)
	v := NewTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// LoadRenderTexture - Load a texture to be used for rendering
func LoadRenderTexture(width, height int32) RenderTexture2D {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ret := C.LoadRenderTexture(cwidth, cheight)
	v := NewRenderTexture2DFromPointer(unsafe.Pointer(&ret))
	return v
}

// UnloadImage - Unload image from CPU memory (RAM)
func UnloadImage(image *Image) {
	cimage := image.cptr()
	C.UnloadImage(*cimage)
}

// UnloadTexture - Unload texture from GPU memory
func UnloadTexture(texture Texture2D) {
	ctexture := texture.cptr()
	C.UnloadTexture(*ctexture)
}

// UnloadRenderTexture - Unload render texture from GPU memory
func UnloadRenderTexture(target RenderTexture2D) {
	ctarget := target.cptr()
	C.UnloadRenderTexture(*ctarget)
}

// GetImageData - Get pixel data from image
func GetImageData(image *Image) unsafe.Pointer {
	cimage := image.cptr()
	ret := C.GetImageData(*cimage)
	return unsafe.Pointer(ret)
}

// GetTextureData - Get pixel data from GPU texture and return an Image
func GetTextureData(texture Texture2D) *Image {
	ctexture := texture.cptr()
	ret := C.GetTextureData(*ctexture)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// UpdateTexture - Update GPU texture with new data
func UpdateTexture(texture Texture2D, pixels unsafe.Pointer) {
	ctexture := texture.cptr()
	cpixels := (unsafe.Pointer)(unsafe.Pointer(pixels))
	C.UpdateTexture(*ctexture, cpixels)
}

// SaveImageAs - Save image to a PNG file
func SaveImageAs(name string, image Image) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cimage := image.cptr()

	C.SaveImageAs(cname, *cimage)
}

// ImageToPOT - Convert image to POT (power-of-two)
func ImageToPOT(image *Image, fillColor Color) {
	cimage := image.cptr()
	cfillColor := fillColor.cptr()
	C.ImageToPOT(cimage, *cfillColor)
}

// ImageFormat - Convert image data to desired format
func ImageFormat(image *Image, newFormat int32) {
	cimage := image.cptr()
	cnewFormat := (C.int)(newFormat)
	C.ImageFormat(cimage, cnewFormat)
}

// ImageAlphaMask - Apply alpha mask to image
func ImageAlphaMask(image, alphaMask *Image) {
	cimage := image.cptr()
	calphaMask := alphaMask.cptr()
	C.ImageAlphaMask(cimage, *calphaMask)
}

// ImageDither - Dither image data to 16bpp or lower (Floyd-Steinberg dithering)
func ImageDither(image *Image, rBpp, gBpp, bBpp, aBpp int32) {
	cimage := image.cptr()
	crBpp := (C.int)(rBpp)
	cgBpp := (C.int)(gBpp)
	cbBpp := (C.int)(bBpp)
	caBpp := (C.int)(aBpp)
	C.ImageDither(cimage, crBpp, cgBpp, cbBpp, caBpp)
}

// ImageCopy - Create an image duplicate (useful for transformations)
func ImageCopy(image *Image) *Image {
	cimage := image.cptr()
	ret := C.ImageCopy(*cimage)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageCrop - Crop an image to a defined rectangle
func ImageCrop(image *Image, crop Rectangle) {
	cimage := image.cptr()
	ccrop := crop.cptr()
	C.ImageCrop(cimage, *ccrop)
}

// ImageResize - Resize an image (bilinear filtering)
func ImageResize(image *Image, newWidth, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResize(cimage, cnewWidth, cnewHeight)
}

// ImageResizeNN - Resize an image (Nearest-Neighbor scaling algorithm)
func ImageResizeNN(image *Image, newWidth, newHeight int32) {
	cimage := image.cptr()
	cnewWidth := (C.int)(newWidth)
	cnewHeight := (C.int)(newHeight)
	C.ImageResizeNN(cimage, cnewWidth, cnewHeight)
}

// ImageText - Create an image from text (default font)
func ImageText(text string, fontSize int32, color Color) *Image {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := color.cptr()
	ret := C.ImageText(ctext, cfontSize, *ccolor)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageTextEx - Create an image from text (custom sprite font)
func ImageTextEx(font SpriteFont, text string, fontSize float32, spacing int32, tint Color) *Image {
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ctint := tint.cptr()
	ret := C.ImageTextEx(*cfont, ctext, cfontSize, cspacing, *ctint)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// ImageDraw - Draw a source image within a destination image
func ImageDraw(dst, src *Image, srcRec, dstRec Rectangle) {
	cdst := dst.cptr()
	csrc := src.cptr()
	csrcRec := srcRec.cptr()
	cdstRec := dstRec.cptr()
	C.ImageDraw(cdst, *csrc, *csrcRec, *cdstRec)
}

// ImageDrawText - Draw text (default font) within an image (destination)
func ImageDrawText(dst *Image, position Vector2, text string, fontSize int32, color Color) {
	cdst := dst.cptr()
	cposition := position.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.int)(fontSize)
	ccolor := color.cptr()
	C.ImageDrawText(cdst, *cposition, ctext, cfontSize, *ccolor)
}

// ImageDrawTextEx - Draw text (custom sprite font) within an image (destination)
func ImageDrawTextEx(dst *Image, position Vector2, font SpriteFont, text string, fontSize float32, spacing int32, color Color) {
	cdst := dst.cptr()
	cposition := position.cptr()
	cfont := font.cptr()
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	cfontSize := (C.float)(fontSize)
	cspacing := (C.int)(spacing)
	ccolor := color.cptr()
	C.ImageDrawTextEx(cdst, *cposition, *cfont, ctext, cfontSize, cspacing, *ccolor)
}

// ImageFlipVertical - Flip image vertically
func ImageFlipVertical(image *Image) {
	cimage := image.cptr()
	C.ImageFlipVertical(cimage)
}

// ImageFlipHorizontal - Flip image horizontally
func ImageFlipHorizontal(image *Image) {
	cimage := image.cptr()
	C.ImageFlipHorizontal(cimage)
}

// ImageColorTint - Modify image color: tint
func ImageColorTint(image *Image, color Color) {
	cimage := image.cptr()
	ccolor := color.cptr()
	C.ImageColorTint(cimage, *ccolor)
}

// ImageColorInvert - Modify image color: invert
func ImageColorInvert(image *Image) {
	cimage := image.cptr()
	C.ImageColorInvert(cimage)
}

// ImageColorGrayscale - Modify image color: grayscale
func ImageColorGrayscale(image *Image) {
	cimage := image.cptr()
	C.ImageColorGrayscale(cimage)
}

// ImageColorContrast - Modify image color: contrast (-100 to 100)
func ImageColorContrast(image *Image, contrast float32) {
	cimage := image.cptr()
	ccontrast := (C.float)(contrast)
	C.ImageColorContrast(cimage, ccontrast)
}

// ImageColorBrightness - Modify image color: brightness (-255 to 255)
func ImageColorBrightness(image *Image, brightness int32) {
	cimage := image.cptr()
	cbrightness := (C.int)(brightness)
	C.ImageColorBrightness(cimage, cbrightness)
}

// GenImageColor - Generate image: plain color
func GenImageColor(width, height int, color Color) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ccolor := color.cptr()

	ret := C.GenImageColor(cwidth, cheight, *ccolor)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientV - Generate image: vertical gradient
func GenImageGradientV(width, height int, top, bottom Color) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctop := top.cptr()
	cbottom := bottom.cptr()

	ret := C.GenImageGradientV(cwidth, cheight, *ctop, *cbottom)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientH - Generate image: horizontal gradient
func GenImageGradientH(width, height int, left, right Color) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cleft := left.cptr()
	cright := right.cptr()

	ret := C.GenImageGradientH(cwidth, cheight, *cleft, *cright)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageGradientRadial - Generate image: radial gradient
func GenImageGradientRadial(width, height int, density float32, inner, outer Color) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cdensity := (C.float)(density)
	cinner := inner.cptr()
	couter := outer.cptr()

	ret := C.GenImageGradientRadial(cwidth, cheight, cdensity, *cinner, *couter)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageChecked - Generate image: checked
func GenImageChecked(width, height, checksX, checksY int, col1, col2 Color) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cchecksX := (C.int)(checksX)
	cchecksY := (C.int)(checksY)
	ccol1 := col1.cptr()
	ccol2 := col2.cptr()

	ret := C.GenImageChecked(cwidth, cheight, cchecksX, cchecksY, *ccol1, *ccol2)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageWhiteNoise - Generate image: white noise
func GenImageWhiteNoise(width, height int, factor float32) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cfactor := (C.float)(factor)

	ret := C.GenImageWhiteNoise(cwidth, cheight, cfactor)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImagePerlinNoise - Generate image: perlin noise
func GenImagePerlinNoise(width, height int, scale float32) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	cscale := (C.float)(scale)

	ret := C.GenImagePerlinNoise(cwidth, cheight, cscale)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenImageCellular - Generate image: cellular algorithm. Bigger tileSize means bigger cells
func GenImageCellular(width, height, tileSize int) *Image {
	cwidth := (C.int)(width)
	cheight := (C.int)(height)
	ctileSize := (C.int)(tileSize)

	ret := C.GenImageCellular(cwidth, cheight, ctileSize)
	v := NewImageFromPointer(unsafe.Pointer(&ret))
	return v
}

// GenTextureMipmaps - Generate GPU mipmaps for a texture
func GenTextureMipmaps(texture *Texture2D) {
	ctexture := texture.cptr()
	C.GenTextureMipmaps(ctexture)
}

// SetTextureFilter - Set texture scaling filter mode
func SetTextureFilter(texture Texture2D, filterMode TextureFilterMode) {
	ctexture := texture.cptr()
	cfilterMode := (C.int)(filterMode)
	C.SetTextureFilter(*ctexture, cfilterMode)
}

// SetTextureWrap - Set texture wrapping mode
func SetTextureWrap(texture Texture2D, wrapMode TextureWrapMode) {
	ctexture := texture.cptr()
	cwrapMode := (C.int)(wrapMode)
	C.SetTextureWrap(*ctexture, cwrapMode)
}

// DrawTexture - Draw a Texture2D
func DrawTexture(texture Texture2D, posX int32, posY int32, tint Color) {
	ctexture := texture.cptr()
	cposX := (C.int)(posX)
	cposY := (C.int)(posY)
	ctint := tint.cptr()
	C.DrawTexture(*ctexture, cposX, cposY, *ctint)
}

// DrawTextureV - Draw a Texture2D with position defined as Vector2
func DrawTextureV(texture Texture2D, position Vector2, tint Color) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	ctint := tint.cptr()
	C.DrawTextureV(*ctexture, *cposition, *ctint)
}

// DrawTextureEx - Draw a Texture2D with extended parameters
func DrawTextureEx(texture Texture2D, position Vector2, rotation, scale float32, tint Color) {
	ctexture := texture.cptr()
	cposition := position.cptr()
	crotation := (C.float)(rotation)
	cscale := (C.float)(scale)
	ctint := tint.cptr()
	C.DrawTextureEx(*ctexture, *cposition, crotation, cscale, *ctint)
}

// DrawTextureRec - Draw a part of a texture defined by a rectangle
func DrawTextureRec(texture Texture2D, sourceRec Rectangle, position Vector2, tint Color) {
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cposition := position.cptr()
	ctint := tint.cptr()
	C.DrawTextureRec(*ctexture, *csourceRec, *cposition, *ctint)
}

// DrawTexturePro - Draw a part of a texture defined by a rectangle with 'pro' parameters
func DrawTexturePro(texture Texture2D, sourceRec, destRec Rectangle, origin Vector2, rotation float32, tint Color) {
	ctexture := texture.cptr()
	csourceRec := sourceRec.cptr()
	cdestRec := destRec.cptr()
	corigin := origin.cptr()
	crotation := (C.float)(rotation)
	ctint := tint.cptr()
	C.DrawTexturePro(*ctexture, *csourceRec, *cdestRec, *corigin, crotation, *ctint)
}
