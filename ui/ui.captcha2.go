package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Captcha2 creates a client-side JavaScript CAPTCHA.
// It returns HTML containing a canvas for the CAPTCHA image, an input field,
// and inline JavaScript to handle the CAPTCHA logic.
//
// IMPORTANT: This is a client-side CAPTCHA and is NOT secure on its own.
// You MUST implement server-side validation to verify the 'js_captcha_verified' field.
func Captcha2() string {
	const captchaLength = 6
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

    b := strings.Builder{}
    b.Grow(captchaLength)
    for i := 0; i < captchaLength; i++ {
        b.WriteByte(chars[rand.Intn(len(chars))])
    }

    captchaText := b.String()

	canvasID := fmt.Sprintf("captchaCanvas_%d", time.Now().UnixNano())
	inputID := fmt.Sprintf("captchaInput_%d", time.Now().UnixNano())
	hiddenFieldID := fmt.Sprintf("captchaVerified_%d", time.Now().UnixNano())

    return Div("", Attr{Style: "display: flex; flex-wrap: wrap; align-items: center; gap: 10px; margin-bottom: 10px; width: 100%;"})(
        // Responsive canvas: width 100% up to a max; height via CSS, real pixel size set in JS
        Canvas("", Attr{ID: canvasID, Style: "border: 1px solid #ccc; width: 100%; max-width: 320px; height: 96px;"})(),
        // Text input becomes full-width on narrow screens
        Input("w-full sm:w-auto flex-1 min-w-0", Attr{ID: inputID, Type: "text", Name: "js_captcha_answer", Placeholder: "Enter text from image", Required: true, Autocomplete: "off"}),
        // Hidden verification field
        Input("", Attr{ID: hiddenFieldID, Type: "hidden", Name: "js_captcha_verified", Value: "false"}),
        Script(fmt.Sprintf(`
        setTimeout(function() {
            const canvas = document.getElementById('%s');
            const ctx = canvas.getContext('2d');
            const input = document.getElementById('%s');
            const hiddenField = document.getElementById('%s');
            const captchaText = '%s';

            function sizeCanvas() {
                const ratio = window.devicePixelRatio || 1;
                const displayWidth = Math.min(320, canvas.clientWidth || 320);
                const displayHeight = 96;
                canvas.width = Math.floor(displayWidth * ratio);
                canvas.height = Math.floor(displayHeight * ratio);
                ctx.setTransform(ratio, 0, 0, ratio, 0, 0);
                canvas.style.width = displayWidth + 'px';
                canvas.style.height = displayHeight + 'px';
            }

            function drawCaptcha() {
                sizeCanvas();
                const w = canvas.clientWidth || 320;
                const h = canvas.clientHeight || 96;
                ctx.clearRect(0, 0, w, h);
                ctx.fillStyle = '#f0f0f0';
                ctx.fillRect(0, 0, w, h);

                ctx.font = 'bold 24px Arial';
                ctx.textBaseline = 'middle';
                ctx.textAlign = 'center';

                for (let i = 0; i < captchaText.length; i++) {
                    const char = captchaText[i];
                    const x = (w / captchaText.length) * i + (w / captchaText.length) / 2;
                    const y = h / 2 + (Math.random() * 10 - 5);

                    ctx.save();
                    ctx.translate(x, y);
                    ctx.rotate((Math.random() * 0.5 - 0.25));
                    ctx.fillStyle = 'rgb(' + Math.floor(Math.random() * 200) + ',' + Math.floor(Math.random() * 200) + ',' + Math.floor(Math.random() * 200) + ')';
                    ctx.fillText(char, 0, 0);
                    ctx.restore();
                }

                for (let i = 0; i < 20; i++) {
                    ctx.beginPath();
                    ctx.arc(Math.random() * w, Math.random() * h, Math.random() * 2, 0, Math.PI * 2);
                    ctx.fillStyle = 'rgba(0,0,0,0.3)';
                    ctx.fill();
                }
            }

            function validateCaptcha() {
				if (input.value.toLowerCase() === captchaText.toLowerCase()) {
					hiddenField.value = 'true';
					input.style.borderColor = 'green';
				} else {
					hiddenField.value = 'false';
					input.style.borderColor = 'red';
				}
			}

            input.addEventListener('input', validateCaptcha);
            drawCaptcha();
            window.addEventListener('resize', drawCaptcha);
        }, 300);
        `, canvasID, inputID, hiddenFieldID, captchaText)),
    )
}
