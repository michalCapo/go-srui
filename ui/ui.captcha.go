package ui

import "html"

func Captcha(siteKey string, secured string) string {
	if siteKey == "" {
		return Div("flex items-center justify-center")(
			Div("")(secured),
		)
	}

	note := Target()
	hidden := Target()
	captcha := Target()

	return Div("")(
		Div("relative flex items-center justify-center")(
			Div("", captcha, Attr{Style: "min-height: 78px; min-width: 304px;"})(),
			// Div("flex-1")(),
			Div("absolute inset-0 flex items-center justify-center opacity-0 pointer-events-none", hidden)(secured),
		),

		Div("text-xs border border-dashed border-black p-1 whitespace-wrap p-2 hidden", note)(
			Div("")("Captcha not loaded, please add the following script to your html file."),
			Div("")(html.EscapeString("<script src=\"https://www.google.com/recaptcha/api.js\" async defer></script>")),
		),

		Script(`
			 setTimeout(function () {
				const note = document.getElementById('`+note.ID+`');
				const captcha = document.getElementById('`+captcha.ID+`');
				const hidden = document.getElementById('`+hidden.ID+`');
				const loaded = window.grecaptcha || null;

				if (loaded == null) {
					setTimeout(function(){
						if (!window.grecaptcha) {
							note.classList.remove('hidden');
						}
					}, 1200);
				} else {
					loaded.ready(function () {
						loaded.render('`+captcha.ID+`', {
							'sitekey': '`+siteKey+`',
							'callback': function () {
								requestAnimationFrame(function(){
									captcha.style.visibility = 'hidden';
									hidden.classList.remove('opacity-0');
									hidden.classList.remove('pointer-events-none');
								});
							},
							'expired-callback': function () {
								requestAnimationFrame(function(){
									captcha.style.visibility = 'visible';
									hidden.classList.add('opacity-0');
									hidden.classList.add('pointer-events-none');
									loaded.reset();
								});
							},
							'error-callback': function () {
								requestAnimationFrame(function(){
									captcha.style.visibility = 'visible';
									hidden.classList.add('opacity-0');
									hidden.classList.add('pointer-events-none');
									loaded.reset();
								});
							},
						});
					});
				}
			}, 300);
		`),

		// Script(`
		// 	window.addEventListener('load', function () {
		// 		const note = document.getElementById('`+note.Id+`');
		// 		const captcha = document.getElementById('`+captcha.Id+`');
		// 		const hidden = document.getElementById('`+hidden.Id+`');
		// 		const loaded = window.grecaptcha || null;

		// 		if (loaded == null) {
		// 			note.classList.remove('hidden');
		// 		} else {
		// 			loaded.ready(function () {
		// 				loaded.render('`+captcha.Id+`', {
		// 					'sitekey': '`+siteKey+`',
		// 					'callback': function () {
		// 						captcha.classList.add('hidden');
		// 						hidden.classList.remove('opacity-25');
		// 						hidden.classList.remove('pointer-events-none');
		// 					},
		// 					'expired-callback': function () {
		// 						captcha.classList.remove('hidden');
		// 						hidden.classList.add('opacity-25');
		// 						hidden.classList.add('pointer-events-none');
		// 					},
		// 					'error-callback': function () {
		// 						captcha.classList.remove('hidden');
		// 						hidden.classList.add('opacity-25');
		// 						hidden.classList.add('pointer-events-none');
		// 					},
		// 				});
		// 			});
		// 		}
		// 	});
		// `),
	)
}
