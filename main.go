// By Maicarons(EldersJavas&gmail.com)
package main

import (
	"flag"
	"fmt"
	"github.com/gen2brain/dlgs"
	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/pkg/browser"
	"github.com/zachomedia/go-bdf"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/language"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var (
	flagTest     = flag.Bool("test", false, "test mode")
	flagEastAsia = flag.Bool("eastasia", false, "East Asia")
	name         = "test"
)

func main() {
	test()
}

//nolint:funlen
func test() error {
	f, _, _ := dlgs.File("BDF file", "*.bdf", false)
	fontBytes, err := ioutil.ReadFile(f)
	if err != nil {
		return fmt.Errorf("failed to load font file: %w", err)
	}
	ubdf, _ := bdf.Parse(fontBytes)
	fmt.Println(ubdf.PixelSize, ubdf.Size, ubdf.Name, ubdf.Encoding)
	nface := ubdf.NewFace()
	println(nface)
	text := `en:      All human beings are born free and equal in dignity and rights.
en-Brai: ⠠⠁⠇⠇⠀⠓⠥⠍⠁⠝⠀⠃⠑⠬⠎⠀⠜⠑⠀⠃⠕⠗⠝⠀⠋⠗⠑⠑⠀⠯⠀⠑⠟⠥⠁⠇⠀⠔⠀⠙⠊⠛⠝⠰⠽⠀⠯⠀⠐⠗⠎⠲
ang:     Ealle fīras sind boren frēo ond geefenlican in ār ond riht.
ar:      يولد جميع الناس أحرارًا متساوين في الكرامة والحقوق.
de:      Alle Menschen sind frei und gleich an Würde und Rechten geboren.
el:      'Ολοι οι άνθρωποι γεννιούνται ελεύθεροι και ίσοι στην αξιοπρέπεια και τα δικαιώματα.
es:      Todos los seres humanos nacen libres e iguales en dignidad y derechos.
eo:      Ĉiuj homoj estas denaske liberaj kaj egalaj laŭ digno kaj rajtoj.
fr:      Tous les êtres humains naissent libres et égaux en dignité et en droits.
got:     ᚨᛚᛚᚨᛁ ᛗᚨᚾᚾᚨ ᚠᚱᛖᛁᚺᚨᛚᛋ ᛃᚨᚺ ᛋᚨᛗᚨᛚᛖᛁᚲᛟ ᛁᚾ ᚹᚨᛁᚱᚦᛁᛞᚨᛁ ᛃᚨᚺ ᚱᚨᛁᚺᛏᛖᛁᛋ ᚹᚨᚢᚱᚦᚨᚾᛋ.
he:      כל בני אדם נולדו בני חורין ושווים בערכם ובזכויותיהם.
hy:      Բոլոր մարդիկ ծնվում են ազատ ու հավասար՝ իրենց արժանապատվությամբ և իրավունքներով:
it:      Tutti gli esseri umani nascono liberi ed eguali in dignità e diritti.
ja:      すべての人間は、生れながらにして自由であり、かつ、尊厳と権利とについて平等である。
ka:      ყველა ადამიანი იბადება თავისუფალი და თანასწორი თავისი ღირსებითა და უფლებებით.
ko:      모든 인간은 태어날 때부터 자유로우며 그 존엄과 권리에 있어 동등하다.
mn:      Хүн бүр төрж мэндлэхэд эрх чөлөөтэй, адилхан нэр төртэй, ижил эрхтэй байдаг.
pl:      Wszyscy ludzie rodzą się wolni i równi pod względem swej godności i swych praw.
pt:      Todos os seres humanos nascem livres e iguais em dignidade e em direitos.
ru:      Все люди рождаются свободными и равными в своем достоинстве и правах.
sw:      Watu wote wamezaliwa huru, hadhi na haki zao ni sawa.
tr:      Bütün insanlar hür, haysiyet ve haklar bakımından eşit doğarlar.
uk:      Всі люди народжуються вільними і рівними у своїй гідності та правах.
vi:      Tất cả mọi người sinh ra đều được tự do và bình đẳng về nhân phẩm và quyền.
zh_Hans: 人人生而自由，在尊严和权利上一律平等。
zh_Hant: 人人生而自由，在尊嚴和權利上一律平等。
`

	if *flagTest {
		text = ""
		for i := 0; i < 256; i++ {
			for j := 0; j < 256; j++ {
				r := rune(i*256 + j)
				if r == '\n' {
					text += " "
					continue
				}
				text += string(r)
			}
			text += "\n"
		}
	}

	path := "example.png"
	if *flagTest {
		if *flagEastAsia {
			path = "test_ea.png"
		} else {
			path = "test.png"
		}
	}
	return outputImageFile(text, *flagTest, path, nface, !*flagTest)
	return nil
}

func glyph(m draw.Image, f font.Face, r byte, off [2]int, x, y, w, h int) {
	draw.Draw(m, image.Rect(x, y, x+w, y+h), &image.Uniform{C: color.White}, image.ZP, draw.Src)
	dot := fixed.P(x+off[0], y+h+off[1])
	dr, mask, mp, _, _ := f.Glyph(dot, rune(r))
	draw.DrawMask(m, dr, &image.Uniform{C: color.Black}, image.ZP, mask, mp, draw.Src)
}

func outputImageFile(text string, grid bool, path string, f font.Face, presentation bool) error {
	const (
		offsetX = 8
		offsetY = 8
	)

	const (
		dotX        = 4
		dotY        = 12
		glyphWidth  = 12
		glyphHeight = 16
	)

	lines := strings.Split(strings.TrimSpace(text), "\n")
	width := 0
	for _, l := range lines {
		w := font.MeasureString(f, l).Ceil()
		if width < w {
			width = w
		}
	}

	width += offsetX * 2
	height := glyphHeight*len(lines) + offsetY*2

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.ZP, draw.Src)
	if grid {
		gray := color.RGBA{R: 0xcc, G: 0xcc, B: 0xcc, A: 0xff}
		for j := 0; j < 256; j++ {
			for i := 0; i < 256; i++ {
				if (i+j)%2 == 0 {
					continue
				}
				x := i*glyphWidth + offsetX
				y := j*glyphHeight + offsetY
				draw.Draw(dst, image.Rect(x, y, x+glyphWidth, y+glyphHeight), image.NewUniform(gray), image.ZP, draw.Src)
			}
		}
	}

	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.Black),
		Face: f,
		Dot:  fixed.P(dotX+offsetX, dotY+offsetY),
	}

	langRe := regexp.MustCompile(`^[a-zA-Z0-9-]+`)

	for _, l := range strings.Split(text, "\n") {
		if presentation {
			if langstr := langRe.FindString(l); langstr != "" {
				lang, err := language.Parse(langstr)
				if err != nil {
					return err
				}
				l = bitmapfont.PresentationForms(l, bitmapfont.DirectionLeftToRight, lang)
			}
		}
		d.Dot.X = fixed.I(dotX + offsetX)
		d.DrawString(l)
		d.Dot.Y += f.Metrics().Height
	}

	fout, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(fout *os.File) {

	}(fout)

	if err := png.Encode(fout, d.Dst); err != nil {
		return err
	}

	return browser.OpenFile(path)

	return nil
}
