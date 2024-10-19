package main

import (
	"bytes"
	//"encoding/base64"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// Define a struct to hold header information
type HeaderInfo struct {
	Key   string
	Value string
}

// Define a struct to hold data to pass to the HTML template
type PageData struct {
	Headers []HeaderInfo
	Image1   string
	Image2   string
	Image3   string
}

// URLs for images
const (
	imageFirecrackerURL 	= "https://s3.nbfc.io/hypervisor-logos/firecracker.png"
	imageQEMUURL        	= "https://s3.nbfc.io/hypervisor-logos/qemu.png"
	imageCLHURL        	= "https://s3.nbfc.io/hypervisor-logos/clh.png"
	imageRSURL        	= "https://s3.nbfc.io/hypervisor-logos/dragonball.png"
	imageContainerURL       = "https://s3.nbfc.io/hypervisor-logos/container.png"
)

// Embed the image data as a base64 string
const imageData2 = `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`
const imageData1 = `iVBORw0KGgoAAAANSUhEUgAAAMUAAABvCAIAAADT3ZMvAAAVyElEQVR42uydaWwU5xnH51OFIkWqVLWV2n4h9PjUD6VH1JRK5TTG3JeB4BobQ4CK0vRDCRImCQaSBmwwBCqDI2wqHIpaH2Abdo13DZgYGxtCGhNjjjhxwAfHnrO7s7M7T/7rgWW9l9cHeLz7/PRovZ6Zd87fvDP7zjvvKxDDjBwCjSkeuLxrrjkm6G0/11t+obcOJyboLX/6RKrqdhOTmD4pRBnNolBHwim7UGHH57CiwibolO+Um6+ZPMQkoE9WWZmgswhVLqHUNDJRbhEqnf+7x1lUQvpkkpSf6q0wYOR8siKjKr/PPiWkT2a38rPn71Nvb09RUfHevfvy8/dHin378gsKDre1tRHDPkX3KT8/f+bM5OnTZ0QPTLNu3Xq3WyaGfYri0/vv/xO6TJ48JXpMmzY9PX2Vw+Ekhn2K4tMHH+yO0afMzNVOp4sY9ol9Yp/YpzEA+8Q+sU/sk1Zhn9gn9ol90irsE/vEPrFPWoV9GrpPGRmZXD7OPg3g044dO/FsDrpEj6SkmStXpnm9CjHsUxSfdDodLmRpaelpaX+OEqtWZRw7dkwUHa2trZ9/3i9UJEmixIN9Clf/yeORRNEeFZvNRkQ5OTuQUU2dOi0wkHvNmJFUWPgRJR6j6VO3S1lxRUy57JjTYEekNLk3XXc4vKPm00NJ+cZJ95xKj0wiDQis86xenRX2Zgs+bd2aTYnHaPp0xeQRzpGgJ0Gn+MJIL1eYvoZQo+GTV6HfGK3jzkjjyk3jzspzPrErCkVHluX16zcgNwr1CZK9++52SjxG06cW+IT3Asos6tGFKD+uNo+WTx6FfnLGLJz1+t500NGrdTb2iX0alk+v+F92qHZPOs8+jTWfmuHTaTHQpx9VmUbRpx9UWQQj+bKoOvpVrZV9GmM+XYVPOHg6HEKPLwz03dPR8ieTWxmvG2GfIHTpPfhEsCezRZzUICNnmnTZ84//O4nYpzHlk11W3rvpXP+p86/XHYgN11377kgUGYj2qtEmnCNfllY+7MBMTjuEM+7KLjcNCfiUlbUW6kyZMjUo+PedVigrLc3Nzd3dnz179hQUFDjt9rIuz8tnnbi/8eVSp4cXlS7MZ3aDHZdR6gN/5YBwRw6X50mSLVu2hJY/IeDTgQMfUh+YWE0leUmJ9+J0bfl09+5dXCySkpJxPIJi1qyU6upqImp57P7311JJp7ukUxpOYCanumSn58kvg9mNrtfO2/943hZLvHZBXNEs3ncqvb29NTU1tf2pqTlnNBotZovkVXAN/f1Fh5rqDxfEGfX2Vms8v96uLZ/a22+pl4+wJYTlFafo+ZB1VRTqSTgjxxq426ujfbddA74g/1KFCRfop6m8Qi0VfhnPVRK05dOtW7dx+Yjk06lTp+n5kNZsF/Q0iBv5MrNQLb/zhZOiYpOVH1aZcZf2NJUFF9miDvaJwNjwqaenp7Oz85vIYOyDBw+oP6taxCH4lNM2FJ+Kv4rn58Rx5ZPBYFy2bPm8efMXLFgYKTB2+fIVtbWGF+CT1a18P8inKuloB/tEYAz4lJubp/7aih6Y5tChf70An0RZ+V6lWah7+oCyhoRaKumM5+ZctOUT7sfhTSSfKga6H9+//0AsVSsxzZEjhS/AJ4DC0rfb5e1tTsQ7bVLuLdcjKZ7LDLTl0+3bdyZPnorjHZqpJCfPKi0tG3M+JRra8gnV+/Py8lCPNqgy5Ouvr9yw4S+wjX3SONrySQVVH8X+YIhafZZ90jha9CkS7JP2YZ/YJ/aJfdIq7BP7xD6xT1qFfRq6T9u/YJ/i2ic8b0ExOoxBCXukUIvaDx48RAHgNUChlnyWlFliinIL/Nt2g32Ka58qKytRjI7C9Cg+YSyaJ9DrayiAA3ckX/2kKgnVPmMJTPxStcPQy42Px7VPiqJ0dHS0DwSmoRAuPXSfvCf/9547pujyXjNBJvrwjmtivTzRYJtosD4Lo31inaM0IbuFiSufYueBRJdMyqVHHjUum8niVmiQeBT65TmrcD5MxU4MnNdgo8QjQX3622cOwUC+tzcrXb5PI6W3iDRIYOBvDVbYE+aG/axnyWU7JR4J6lNGiyjoFNyAPwk9LbhsH4JPvzNG9GlpI/uUMKy5KuKQPzv8OlrSyD6xT+yTxkhQn9Ku2AUj4ebpSdTRnAb2iX0aKrtvOscbZTSo8spZCz7Hn/dsbXUOwadfGyyRfFrE9+PxxOPHj0tKSo4ePVoUAP4tLi6+ceOGV6Eup7fbpaiB77JCg0VW86cw5QUe4QItb2Kf4gh02Kp2sxnULi8Gbty40SPLNBLoe+R5TS40Zjc3IFIu2VNbpE/NXko84tanbdveDvtsGEq98cY6p5MfvbFPg2H79pxIPqHRJpcrnl/6Dg/7xD6NOdgnhn1in7QK+8SwTzGA5itRNBC2JcKMjNX8+459GhxohnDp0tSgxnrmz1+wePFSVPb1ehOxcIh9GhYPHz7s8dHrj+5u37/cyRj7xIwN2CeGfWK0CvvEsE+MVmGfGPaJ0SrsE8M+MVqFfWLYJ0arsE8M+8RoFcHfFS568fI9hO8D3U5arVYKwOFwYGDPQGAaURSpD7fbHWMSi8VCfXg8HlQKwLABk+D1OhoSZrO5sbHx+PHjeXl70Z4dmkg8ceI/LS0tQdsLom++uouiLCV0epPJROFA5RmMDd1qSZKit/ofugi7XaRwYFV9Y/vAAcKaBC4Oa4sKhiPp0927X6JuECqg+TtXRae9/bsCM2AsembC2CiBhP6m327evDlnztwBk6SkzD58+Ij/MGRkZKCNuQGXgo57aZCg57uDBw8uWbIUrdShlqZ/bviOIejSo7j4mCg6ItX2xEID1wEr+eabf5dlmcLx1ltbQqZPRschYRVsaGgI2rf4jiH19fUUmU2bNgXtKCwxJ2cHhaOg4DAm7uu8JC09fdXChYsCF4d9snbt2uzsbSdPnrx//36Q65APZ29XVxfeg21tbTWZzKiN2NHxFT79zbg9evSou7v7mU/oaHDu3HnqnkVgGUFr1tTUqNZvxNgogbU0Guv8xy85OQXzHDBJYeFH1AfOkszM1WqXr5EDqzcTB4wGA9YqNXU5tgubENp/ld+w69c/C1uPCnXxgtYK0+NMwDZSCMiY8Yof5hbUziKiubmZQsjP34+dEDQx5o9zOEpHy1ifoH2LIampyyyWMMp+/PGJoJ7cQheH5NOnJy1atBgZNhR5mim0ozsdnOTYWPQtiO86nW7XrvdgC2osUh/YKGRGmzdvfuYTHMSqYL7+FnB37tzV36cmLDjwSKiL7x8zsNIXL9b7fZo9ey7mGVh3OzQJzpuiomLqAzl8VtYaTOZPguRhk2RnZ1PMYBeouyxQIMxHPSTqp7qL279t79x/qziuOP7H0UJKHAitC9hGJsQODdgEEz9A1JJLSYXsguIqLqY8/ANxJMK7EJuKuMKY2HXlFIxRYsRD1Cb+LaRJ81NDo83HPvJXxzPLaq9vhFC0R+urq9k5szNnvnNeM757717MPjw8LN5AArx9OtUSvfbaFrHoAmH8z7uqaciAL6hsc8xzk2fQ2bPneLoq6wvX5ORkXP/MmbNWX5cJZAGRPwvK6WdPz5+NESU0NzeHP8AtLA9LC4vZ3NyC+uS981anu7ubeUf5JcIT9VB6Hk+Bfrpz547HEzXR9kyqp127djU27mDpGAsO2auv/lJtwrtt27aAhZ41NDSMj4/LjeM9Un7yams3ByxtbW085fz580k+mp6+W1Hxip8w+97U1IQgjh492tX1xzff3GZrPVXfdHV12WRQgeWLuCUlnLC4PhJHy8YQpOTtt5tjTWOaLAbruXPPHOPu3XuEj8rKX9Py4iNWYrXj+qdPn6G+bxz7jmwBgQlH0ySoDQ9fF/vs7Cwl6CHzcffu3QsjY0QTsx6YDvQWZjQBTxJBdXXNihU/10p6990lkgKkwhOfeAOBoZVD7b/zZlU6pwU3NjaWah1iSWkC+vv7s1myiT60trZag2oTQQjBUhKID6QyqNiTRc8jGVNp/f0fACm+q6nYhcLhlZZl4QJWq48EmDz8X1/5woULNGuNV1VV8UVT8OGHp5M0wpXBVTKFiofKz36gsI2RRlpaWuPT8TQlIcC4adMmRkoYhMMKVvhdZOAlUFo7dXX1yFmg5y5CM+3L6PBE8QHwou7evVtfv7W7+09IaUm+oL6+Hh4Nhp8isXI5WB5PuIsPHjxIMglB8wyPJzkEGdTWttvjqa+vLymDWFImdzW4YcMG1kZQTd5bHFKNjf1jQSzz/cclAA0yT+b8RSYSxfylDZwKvK6Yd8tiUySE69ev+x8kbm//LeO1mmhKKmgKwG6SRrynFBYbDmAlcJOvwieuDLFbwDI09HePp82bN2ukck5Y/F6ncn366b88noiLbYFhzfFqABw9OXnyJPr7o48G8PGX4Kmurs7j6cSJEz8BPGHRGIt3ulmLSSn03ns91gK96ujooIRZVw/5Er/VGLzirhqeUEi8+xqsmBxo6siRv1BHbuv69evNQea/ArFK9FAt46cnadTR8TvumnUj30EJmtXDZXQ0tAOffDLqK9TW1mKRIpj+lad7BUmc5PFEHGpuDP84REKgp6fnrXnayZLjgrEsPD169Cg/nuRgloonqcmSSGYRM89zJUcEgTEq6X2hNTU1MNoEX7x4kUJcN4+nfft+H9sjbAFcFikzGW+8sZWxWH0mQNHTxMQEqkuuGA4v9bPxxHRi42hZ0LF0gBPaSnnTItRJBp7kQ9syUDv8Ro3wRDnv27XVwl0s3eTkbXqLWkKkQ0NDqOqy8ITh/F8aZeDp6tWrcX2sjOQb4+n48eOpLN5FyAiqCUOcsVPQmpfQ8Oo/PoqFGlNTU1rH3GWM3367ZG5mZmZ4rkJxcjOHDh3i6dYOgTcWymqSeJPlQnmPjNzweMIpSSIaGBiksj26srISEFCIDZVZhxEjyFoqFU/Mgg+G+EKoRKGtENyyGzduGBKYI1BFjhdnF0eKQgJA4rPl4EkPI6wgJPQXBhhbkIonmfaABcNPVKjBB3ha8CrWBizoGLSOHOcMunnzpibeBN3XV5q2w5ogChvv1q2/AcTmkFosLBuKmvFcrDTd5cKosZDIdGhdjYyMWKqQyIheWSPkA8k+8Div+TJ+14pPJeGY3dWrVwviL71k0M+LJ9E77/zBV2OYLF2DGr3lM/r+vbBI6LNMPOl5wYXGI48c4ymDhQcRLNCVGE+CVMDC3YqKCsvGZhPTQ33vEBBMJbkJURIGy06RWfDp74zYhY0BecesOtqZmfk3wrGeACzMk7nA6CobFHGi5ck8nmRuRCg2mSQ6gAGSqC1KEC9u4jLwRKLSV8Nqo4SS3FQWnqgc/5QgAM/GU8yCbwHeAzxlsJAcB0+EMDnTmB5P5gDlJLSOMsjwTk9/7tKDZ4CF5LBzZ5NnJC0iPKGz0WcsmC1bXqemmd2mpl22hcVIFzq2kp1ESq5cuZKpn6jwN+XEwaJfVPwTvedlYZevnzAFeBfPA08oDYJG7LS/sL4EkBl4IsUVsDQ0NPb2HhFLbO9wEeKnICzWTR7vJwhY+DHWJDfhDstI0Q4g2L1IrAGf2iUbhP31wTn5KnlX5Hh8VGhJI6aT9mXsPv54iDpE3cIEhXgkZlZEbLGpEZZic3Or9WfPnj2kE61Lei6PKBVPxApL/adGdSAPleWP85amUuM73MZS4zttUCyD7t+/j5VxPh8rfl9OXvCKsbPOi52O2RVs3aBCCc0CLaIlbkchiICUteIiiYDzZIlpNI2lxAYHr3g8sUeLDL2xU95BekiXyrUA8JFLwhP9VPteQT4nPOHxlYqna9euPc98AZhAI7oNn3mH4MmTLzOZtF85ofEKBMHlN+Y6O7v8Zhk9X0ygv24zNzv7GDUmaaDtiLT54s3l6OgoYlLyffv27fIszRpaszm7hOhKwhP7u5JVSepc9NPPZ3Z2dgabLceOHcvDyGanZyRZEFw0xag0Q4ScOr/BBojwBG6ePv2/JcMaGxu9KVeGkzG6TXdohXZF/LEk7Fp2l4Ktd4xXjCflx4UnTRa7q0F+3DRcYe/8fsuEXCi5HeasxMSICM1sS0E7GHQGPYGtebKU0BZU8NKfmrqzmGWYd4w0qd9999TKe3t7rTw45sCmkBw+jydCOWVf5+a+wCyqS+3t7U8iYn8Go+y9OkLIVDyxWxdkCvFigy1zfMRSnfHs/bu+jP1g8HT79m0il2/m6b9fffWfx4+/oPdMCfuLz9oPHhgYYFnAYDuRtEllsn+YTiqn7gczTrFAJNZg4UEPHz7Mk+nGnSS2ZzhBwMj7qGnBZyBPnTpFeowpUYAmN4KcahIRDjincXw1znxSLtc7xhNpJ69ClFAgQZWKJ/qjsOPSpUveU04NVD/77HNz8OVCkfcK8KRFRXYb+bMGCKFw/F0gbLzkV68t83ym3E+Pp8OHDwc53+CoDckuYF5VVVNdvQk9xODx5rhQy9q0Z/PVSxD3hafAwgXU1q37FfVhRP0CTeMiD2sjVygkFi729lm19qx79+4nOQj8kQNVm9IKLHccXgJGEmBkTSmx80wChKZ8fPyfSURxyqe1tc2iIRAc4ElbJUggsEpoGrUpPCkjj9gp1w6dVmbq2JGhFzj1MfcBnnQhRuRP5Ti8kDtYMuU//0SuOc5nWngC+fwkI9c21tq16+L8U5yfrKmpRgmZk0EaOhhezGJnKdlAzJ9JQnwSaJzZUrqIfU36wPJQCV4zpiRJo4MHD2Jf1M81a9bY7geHulLxBNpaWuxW+iEn4UlaBCVqxoF1pTQmmZpgOyVOtMqu4dUF+8G6K53kRc2IDhw4gC4oC0860qpMdxArcvZFoo8vjydyJHK5Vq36hdrMYCHXZ44nn+RRKMlmoU2cAyxskpump6c5QGdKSAvDi9JGTVjHqW11gBJOwKXmYCxdRAUPcTuDxnCsBT7JlhmetGHnWXguHfNCppRy3eUYI+UcAhMXIGO+k2cQ7iYV1D6ft27dohzjpRbixWmdt+NcSqyXhSdyD2hvkiWYcy4WKO/TCc5n2vYZdzMuGBWm4sMiWbWZwUJaljkzI4L/hC+czUKFHTve0kHpnMSaGxwcxFd4+eVXFkzeKkRscuRiaCSpv/76G86TLHZg/caNG7Wnkeq/owBg1EDef7+fcgyftcAnq8viO8j2khEILHa3ubmZJeSNAHuidpdP/nBSKd+/fz+H360QixaHNa6FW3RDLfCsy5cv21YB5V6GqkBYh1bDT2W7VyFquXhiOhkYRIsQXwKNypYIxfNXJsEIJvwUqs2YUp/F9zwstGwQLJXgwilhJ5//5eifpw/IIjLN+p8tycEGm2SSBqheqTAemtrnjmQVCNla0icllFtlfSaZFM+jxUaxVK0mjyj+n7OgF5QKPBVU4KmgF5UKPBVU4KmgF5UKPBVU4KmgF5UKPBX0Y9IPv5AVmfNXy6cAAAAASUVORK5CYII=`



// HTML template to generate the response
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Serverless Demo</title>
        <style>
        table {
            border-collapse: collapse;
            width: 50%;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
</head>
<body>
<h1> Hello <img src="data:image/png;base64,{{.Image1}}" alt="CAMAD Logo" height="60px"/>
    </h1>
    <h2> RuntimeClass
    </h2>
    <img src="{{.Image3}}" alt="Runtime Class" height="100px" />
    <h1>Request Headers</h1>
    <table>
        <tr>
            <th>Header</th>
            <th>Value</th>
        </tr>
        {{range .Headers}}
        <tr>
            <td>{{.Key}}</td>
            <td>{{ .Value | safeHTML }}</td>
        </tr>
        {{end}}
    </table>

    <!--<h1>Request Headers</h1>
    <ul>
    {{range .Headers}}
        <li><strong>{{.Key}}:</strong> {{.Value}}</li>
    {{end}}
    </ul>-->
    <h2> Brought to you by 
    </h2>
    <img src="data:image/png;base64,{{.Image2}}" alt="Nubificus LTD"/>
</body>
</html>
`

func safeHTML(s string) template.HTML {
    return template.HTML(s)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("On my way to unikernels!!!")
	fmt.Println("Received request")

	// Create header info for the template
	headers := make([]HeaderInfo, 0, len(r.Header))
	for key, values := range r.Header {
		for _, value := range values {
			headers = append(headers, HeaderInfo{
				Key:   html.EscapeString(key),
				Value: html.EscapeString(value),
			})
		}
	}

	// Determine the runtime class based on the hostname
	host := r.Host
	var imageURL string
	if strings.Contains(host, "hellofc") {
		imageURL = imageFirecrackerURL
	} else if strings.Contains(host, "helloqemu") {
		imageURL = imageQEMUURL
	} else if strings.Contains(host, "helloclh") {
		imageURL = imageCLHURL
	} else if strings.Contains(host, "hellors") {
		imageURL = imageRSURL
	} else {
		imageURL = imageContainerURL // Default image or error image if needed
	}


	data := PageData{
		Headers: headers,
		Image1:   imageData1,
		Image2:   imageData2,
		Image3:   imageURL,
	}

	// Parse and execute the template
	tmpl, err := template.New("page").Funcs(template.FuncMap{"safeHTML": safeHTML}).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the content type to HTML and write the output
	w.Header().Set("Content-Type", "text/html")
	buf.WriteTo(w)
}

func main() {
	log.Print("helloworld: starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("helloworld: listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
