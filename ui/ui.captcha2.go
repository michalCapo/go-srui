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
	for range captchaLength {
		b.WriteByte(chars[rand.Intn(len(chars))])
	}

	// captchaText := b.String()
	captchaText := "Hello"

	canvasID := fmt.Sprintf("captchaCanvas_%d", time.Now().UnixNano())
	inputID := fmt.Sprintf("captchaInput_%d", time.Now().UnixNano())
	hiddenFieldID := fmt.Sprintf("captchaVerified_%d", time.Now().UnixNano())

	return Div("", Attr{Style: "display: flex; align-items: center; gap: 10px; margin-bottom: 10px;"})(
		Canvas("", Attr{ID: canvasID, Width: 250, Height: 100, Style: "border: 1px solid #ccc;"})(),
		Input("", Attr{ID: inputID, Name: "js_captcha_answer", Placeholder: "Enter text from image", Required: true, Autocomplete: "off"}),
		Input("", Attr{ID: hiddenFieldID, Name: "js_captcha_verified", Value: "false"}),
		Script(fmt.Sprintf(`
		setTimeout(function() {
			const canvas = document.getElementById('%s');
			const ctx = canvas.getContext('2d');
			const input = document.getElementById('%s');
			const hiddenField = document.getElementById('%s');
			const captchaText = '%s';

			function drawCaptcha() {
				ctx.clearRect(0, 0, canvas.width, canvas.height);
				ctx.fillStyle = '#f0f0f0';
				ctx.fillRect(0, 0, canvas.width, canvas.height);

				ctx.font = 'bold 24px Arial';
				ctx.textBaseline = 'middle';
				ctx.textAlign = 'center';

				for (let i = 0; i < captchaText.length; i++) {
					const char = captchaText[i];
					const x = (canvas.width / captchaText.length) * i + (canvas.width / captchaText.length) / 2;
					const y = canvas.height / 2 + (Math.random() * 10 - 5);

					ctx.save();
					ctx.translate(x, y);
					ctx.rotate((Math.random() * 0.5 - 0.25));
					ctx.fillStyle = 'rgb(' + Math.floor(Math.random() * 200) + ',' + Math.floor(Math.random() * 200) + ',' + Math.floor(Math.random() * 200) + ')';
					ctx.fillText(char, 0, 0);
					ctx.restore();
				}

				for (let i = 0; i < 20; i++) {
					ctx.beginPath();
					ctx.arc(Math.random() * canvas.width, Math.random() * canvas.height, Math.random() * 2, 0, Math.PI * 2);
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
		}, 300);
	`, canvasID, inputID, hiddenFieldID, captchaText)),
	)
}
