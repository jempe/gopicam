/*--include:parts/common.js:--*/
const requestInit = {
	"cache": "no-store",
	"mode": "same-origin",
	"credentials": "same-origin",
	headers: {
		"Content-Type": "application/json; charset=utf-8"
	}
};

function handleResponse(response)
{
	if (document.querySelectorAll('button.disabled').length > 0)
	{
		document.querySelector("button.disabled").classList.remove("disabled");
	}

	if (response.status >= 200 && response.status < 300)
	{
		return Promise.resolve(response)
	} 
	else
	{
		if(response.status == 401)
		{
			show_login();
		}

		return Promise.reject(new Error(response.statusText))
	}
}
function handleJson(response)
{
	return jsonResponse = response.json();
}

// display error messages
function log_error(message)
{
	console.log(message);
}

// get data of form
function get_form_data(form)
{
	let form_data = {};

	let fields = form.querySelectorAll("input, select");

	for(let field_index in fields)
	{
		let field = fields[field_index];

		form_data[field.name] = field.value;
	}

	return form_data;
}

// get data of form
function populate_form_data(form, form_data, csrf_token)
{
	let fields = form.querySelectorAll("input, select");

	for(let field_index in fields)
	{
		let field = fields[field_index];

		if(field.name == "gorilla.csrf.Token")
		{
			field.value = csrf_token;

		}
		else if(form_data[field.name])
		{
			field.value = form_data[field.name];
		}
	}
}

/*--includeend--*/

// show login form
function show_login()
{
	document.body.classList.add("show_login");
}

// hide login form
function hide_login()
{
	document.body.classList.remove("show_login");
}


//submit login form
function submit_login_form(form)
{
	document.querySelector(".login_form_container").classList.remove("error");

	let loginRequest = { 
		"cache": "no-store",
		headers: {
			"Content-Type": "application/x-www-form-urlencoded"
		},
		"method" : "POST",
		"body" : "username=" + encodeURIComponent(document.getElementById("username").value) + "&password=" + encodeURIComponent(document.getElementById("password").value)
	}

	fetch("/api/login" , loginRequest).then(handleResponse).then(handleJson).then(function(data)
	{
		if(data.access != "granted")
		{
			//if can't login
			if(document.querySelectorAll("#login_form span.error").length == 0)
			{
				document.getElementById("login_form").insertAdjacentHTML('afterBegin', '<span class="error">Wrong username or password</span>');
			}

			document.querySelector(".login_form_container").classList.add("error");
		}
		else
		{
			hide_login();
			get_preview();
		}

	}).catch(function(error)
	{
		log_error('Request failed' +  error);
	});

	return false;
}

// preview image and camera status
function get_preview()
{
	// prepare request
	let previewRequest = requestInit;
	previewRequest["method"] = "GET";

	// get data from preview api
	fetch("/api/camera/preview", previewRequest).then(handleResponse).then(handleJson).then(function(data)
	{
		document.querySelector("main").dataset.status = data.status;
		console.log(data.status);

		// create image
		let preview_image = new Image();

		preview_image.addEventListener('load', function() {

			let preview_canvas;

			if(document.getElementById("preview") == null)
			{
				// create canvas element if it doesn't exist
				preview_canvas = document.createElement("canvas");
				preview_canvas.id = "preview";
				preview_canvas.width = preview_image.width;
				preview_canvas.height = preview_image.height;

				document.getElementById("preview_container").appendChild(preview_canvas);
			}
			else
			{
				preview_canvas = document.getElementById("preview");	
			}

			let preview_ctx = preview_canvas.getContext('2d');

			preview_ctx.drawImage(preview_image, 0, 0);
		}, false);

		preview_image.src = data.image;

		document.getElementById("background").style.backgroundImage = "url(" + data.image + ")";

		// Camera status values: md_video, md_ready, ready, video, halted, tl_md_ready 

		setTimeout(get_preview, 1000);
	}).catch(function(error)
	{
		log_error('Request failed' +  error);
	});
}

// send command to camera API
function send_command(camera_command) {
	// prepare request
	let commandRequest = requestInit;
	commandRequest["method"] = "GET";

	let command_url = "";

	if(camera_command == "power")
	{
		if(document.querySelector("main").dataset.status == "halted")
		{
			command_url = "/api/camera/start"; 
		}
		else
		{
			command_url = "/api/camera/stop";
		}
	}
	else if(camera_command == "motion")
	{
		if(/^md_(video|ready)$/.test(document.querySelector("main").dataset.status))
		{
			command_url = "/api/camera/motion_detect/stop"; 
		}
		else
		{
			command_url = "/api/camera/motion_detect/start";
		}
	}
	else if(camera_command == "timelapse")
	{
		if(/^(tl_md_video|timelapse|tl_md_ready)$/.test(document.querySelector("main").dataset.status))
		{
			command_url = "/api/camera/timelapse/stop"; 
		}
		else
		{
			command_url = "/api/camera/timelapse/start";
		}
	}
	else if(camera_command == "record")
	{
		if(/^video$/.test(document.querySelector("main").dataset.status))
		{
			command_url = "/api/camera/record/stop"; 
		}
		else
		{
			command_url = "/api/camera/record/start";
		}
	}
	else if(camera_command == "photo")
	{
		command_url = "/api/camera/photo/take"; 
	}

	if(command_url != "")
	{
		// send command to api
		fetch(command_url, commandRequest).then(handleResponse).then(handleJson).then(function(data)
		{
			console.log(data);
		}).catch(function(error)
		{
			log_error('Request failed' +  error);
		});
	}
}

function toggleFullScreen()
{
	if( ! document.fullscreenElement)
	{
		document.body.requestFullscreen();
	} 
	else
	{
		if(document.exitFullscreen)
		{
			document.exitFullscreen(); 
		}
	}
}

get_preview();
