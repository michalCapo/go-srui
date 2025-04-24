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
		Div("flex items-center justify-center")(
			Div("", captcha)(),
			// Div("flex-1")(),
			Div("pointer-events-none opacity-25 hidden", hidden)(secured),
		),

		Div("text-xs border border-dashed border-black p-1 whitespace-wrap p-2 hidden", note)(
			Div("")("Captcha not loaded, please add the following script to your html file."),
			Div("")(html.EscapeString("<script src=\"https://www.google.com/recaptcha/api.js\" async defer></script>")),
		),

		Script(`
			 setTimeout(function () {
				const note = document.getElementById('`+note.Id+`');
				const loaded = window.grecaptcha || null;

				if (loaded == null) {
					note.classList.remove('hidden');
				} else {
					loaded.ready(function () {
						loaded.render('`+captcha.Id+`', {
							'sitekey': '`+siteKey+`',
							'callback': function () {
                                const captcha = document.getElementById('`+captcha.Id+`');
                                const hidden = document.getElementById('`+hidden.Id+`');

								captcha.classList.add('hidden');
								hidden.classList.remove('hidden');
								hidden.classList.remove('opacity-25');
								hidden.classList.remove('pointer-events-none');
							},
							'expired-callback': function () {
                                const captcha = document.getElementById('`+captcha.Id+`');
                                const hidden = document.getElementById('`+hidden.Id+`');

								captcha.classList.remove('hidden');
								hidden.classList.add('hidden');
								hidden.classList.add('opacity-25');
								hidden.classList.add('pointer-events-none');
                                loaded.reset();
							},
							'error-callback': function () {
                                const captcha = document.getElementById('`+captcha.Id+`');
                                const hidden = document.getElementById('`+hidden.Id+`');

								captcha.classList.remove('hidden');
								hidden.classList.add('hidden');
								hidden.classList.add('opacity-25');
								hidden.classList.add('pointer-events-none');
                                loaded.reset();
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
