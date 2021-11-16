
gf_picker__main();

//---------------------------------------------------
function gf_picker__main() {

    console.log("gf_page_picker");

    const api_host_str = "gloflow.com"

    // import jquery if its not defined.
    // testing both condition, because on some sites window.jQuery is defined 
    // but $ is not defined and vice-versa.
    if (!window.jQuery || typeof $ === 'undefined') {
        console.log("GF - jquery not defined - inserting");

        let s = document.createElement("script");
        s.setAttribute('crossorigin', 'anonymous');
        s.setAttribute('integrity',   'sha256-cCueBR6CsyA4/9szpPfrX3s49M9vUU5BgtiJj06wt/s=');
        s.setAttribute('src',         'https://code.jquery.com/jquery-3.1.0.min.js');
        s.setAttribute('type',        'text/javascript');
        document.body.appendChild(s);
        
        s.onload = () => {
            gf_picker__create_ui(api_host_str);
        }
    } else {
        gf_picker__create_ui(api_host_str);
    }
}

//---------------------------------------------------
function gf_picker__create_ui(p_api_host_str) {

    // CSS
    $("body").append(`
    <style>
    
    div#gf_page_picker {
        
        position: absolute;
        width: 100%;
        height: 100%;
        top: 0px;
        left: 0px;
        
        /*important so that the GF UI is above all other page content*/
        z-index: 100000;

        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
	    margin:      0;

        color: black;
        font-size: 12px;
    }

    div#background {
        position: absolute;
        width: 100%;
        height: 100%;
        background-color: gray;
        opacity: 80%;
        top: 0px;
        left: 0px;
    }

    div#gf_bookmark {
        top:              30px;
        left:             19%;
        position:         relative;
        background-color: yellow;
        width:            800px;
        padding:          10px;
        padding-bottom:   7px;
    }

    div#gf_bookmark input {
        border-width: 0px;
    }

    div#gf_bookmark div#description {
        width: 100%;
    }
    div#gf_bookmark div#description input {
        height:       60px;
        width:        100%;
        padding:      0px;
        /*border-width: 2px;*/
        border-color: #808080;
        border-style: solid;

        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
        font-size:   16px;
    }


    div#gf_bookmark div#tags {
        width: 100%;
    }
    div#gf_bookmark div#tags input {
        width: 100%;
        padding: 0px;

        font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
        font-size:   14px;
    }



    div#gf_bookmark div#close_btn {
        position:         absolute;
        top:              0px;
        right:            -49px;
        width:            50px;
        height:           50px;
        background-color: #fcfd30;
        text-align:       center;
        color:            #7d5d20;

        cursor: pointer;
    }

    div#gf_bookmark div#close_btn:hover {
        opacity: 0.9;
    }

    div#gf_bookmark div#close_btn img {
        width: 100%;
        position: absolute;
        top: 0px;
        left: 0px;
    }

    div#gf_bookmark div#submit_btn {
        background-color: gray;
        text-align:       center;
        padding-top:      11px;
        padding-bottom:   10px;
        width:            100%;
        cursor:           pointer;
        color:            white;
        opacity:          1;
    }

    </style>`);



    const page_picker_element = $(`
        <div id="gf_page_picker">
            <div id="background"></div>
        </div>`);

    // PAGE_PICKER
    $("body").append(page_picker_element);

    


    // BOOKMARK
    const current_url_str  = window.location.href;
    const bookmark_element = $(`
        <div id="gf_bookmark">
            <div id="url">${current_url_str}</div>
            <div id="description">
                <input value="url description"></input>
            </div>
            <div id="tags">
                <input value="tags"></input>
            </div>
            <div id='close_btn'>
                <img src='https://gloflow.com/images/static/assets/gf_close_btn_small.svg'></img>
            </div>
            <div id='submit_btn'>ok</div> 
        </div>`);

    $("#gf_page_picker").append(bookmark_element);

    // SUBMIT_BTN
    $("div#gf_bookmark div#submit_btn").on('click', function(p_event) {
        const submit_btn = p_event.target;

        const url_str         = current_url_str;
        const description_str = $(bookmark_element).find("#description input").val();
        const tags_lst        = $(bookmark_element).find("#tags input").val().split(" ");
        gf_picker__create_bookmark__http(url_str,
            description_str,
            tags_lst,
            p_api_host_str,
            // on_complete
            ()=>{
                $(submit_btn).css("background-color", "green");
            },
            // on_error
            ()=>{
                $(submit_btn).css("background-color", "red");
            });
    })

    //-------------------------------
    // CLOSE - via close_btn
    $("div#gf_bookmark div#close_btn").on('click', function() {
        $(page_picker_element).remove();
    })

    // CLOSE - via background click
    $(page_picker_element).find("#background").on('click', function() {
        $(page_picker_element).remove();
    });

    //-------------------------------
}

//---------------------------------------------------
function gf_picker__create_bookmark__http(p_url_str,
    p_description_str,
    p_tags_lst,
    p_api_host_str,
    p_on_complete_fun,
    p_on_error_fun) {
        
    const url_str = `https://${p_api_host_str}/v1/bookmarks/create`
    const data_map = {
        "url_str":         p_url_str,
        "description_str": p_description_str,
        "tags_lst":        p_tags_lst,
    };

	$.post(url_str,
		JSON.stringify(data_map),
		()=>{
            p_on_complete_fun();
        },
        "json")
        .fail(()=>{
            p_on_error_fun();
        });
}

//---------------------------------------------------
/*function gf_picker__create_screenshot() {
    const capture = async () => {
        const canvas  = document.createElement("canvas");
        const context = canvas.getContext("2d");
        const video   = document.createElement("video");

        try {
            const captureStream = await navigator.mediaDevices.getDisplayMedia();
            video.srcObject     = captureStream;


            context.drawImage(video, 0, 0, window.screen.width, window.screen.height);

            console.log("drawn")
            console.log(window.screen.width)
            console.log(window.screen.height)
            const frame = canvas.toDataURL("image/png");
            captureStream.getTracks().forEach(track => track.stop());
            


            // const canvas_viewer = document.createElement("canvas");
            const canvas_viewer = $(`<canvas id='mycanvas' width='${window.screen.width}' height='${window.screen.height}'></canvas>`);
            $("#gf_page_picker").append(canvas_viewer);
            var myImage = new Image();

            console.log(frame)
            myImage.src = frame;
            
            console.log(canvas_viewer.get())
            canvas_viewer.get()[0].getContext("2d").drawImage(myImage, 0, 0, 400, 400);

            window.location.href = frame;
        } catch (err) {
            console.error("Error: " + err);
        }
    };

    capture();
}*/