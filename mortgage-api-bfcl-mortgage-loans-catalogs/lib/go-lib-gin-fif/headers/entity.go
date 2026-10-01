package headers

var HeaderKeys = []string{
	"X-Forwarded-Proto",
	"X-Forwarded-For",
	"X-Channel",
	"X-Commerce",
	"X-Forwarded-Host",
	"X-Session-Id",
	"X-Country",
	"X-Forwarded-Port",
	"X-Functionality",
	"X-True-Client-Ip",
	"X-Trace-Id",
	"X-Date-Time",
	"X-App-Version",
	"User-Agent",
	"X-User-Agent",
	"X-App-Platform",
	"X-Customer-Hash",
	"X-Request-Id",
	"x-Authorization",
}

type MandatoryHeaders struct {
	XChannel        string `header:"X-Channel" binding:"required,oneof=Web Mobile Kiosco OM2 OM1 PCOM CAJA CYBER WEB_PORTAL BACKOFFICE PUBLIC_SITE LIA SALONES_VIP OMNITELLER SMARTIX INBROKER SFWEB PF Portal IVR"`
	XCommerce       string `header:"X-Commerce" binding:"required,oneof=BANCO CMR FALABELLA SODIMAC DICICO TOTTUS SAGAFALABELLA LINIO GLOBALSELLERCENTER REALSTATE SEGUROS FPAYPAYMENTS LOYALTY BANCAEMPRESA IVR XERPA"`
	XCountry        string `header:"X-Country" binding:"required,oneof=AR CL CO PE UY MX"`
	XForwardedProto string `header:"X-Forwarded-Proto" binding:"required"`
	XForwardedFor   string `header:"X-Forwarded-For" binding:"required"`
	XForwardedHost  string `header:"X-Forwarded-Host" binding:"required"`
	XSessionId      string `header:"X-Session-Id" binding:"required"`
	XForwardedPort  string `header:"X-Forwarded-Port" binding:"required"`
	XFunctionality  string `header:"X-Functionality" binding:"required"`
	XTrueClientIp   string `header:"X-True-Client-Ip" binding:"required"`
	XTraceId        string `header:"X-Trace-Id" binding:"required"`
	XDateTime       string `header:"X-Date-Time" binding:"required"`
	XAppVersion     string `header:"X-App-Version" binding:"required"`
	XRequestId      string `header:"X-Request-Id" binding:"required"`
	UserAgent       string `header:"User-Agent" binding:"required"`
	XUserAgent      string `header:"X-User-Agent"`
	XAppPlatform    string `header:"X-App-Platform" binding:"required"`
	XCustomerHash   string `header:"X-Customer-Hash" binding:"required"`
	XAuthorization  string `header:"X-Authorization"`
}
