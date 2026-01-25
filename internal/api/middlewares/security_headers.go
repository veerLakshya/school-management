package middlewares

import (
	"net/http"
)

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off") // disables dns prefetching requres

		w.Header().Set("X-Frame-Options", "DENY")                                                 // prevents embedding of page in iframe
		w.Header().Set("X-XSS-Protection", "1;mode=block")                                        // blocks page if xss attack is detected (An XSS (Cross-Site Scripting) attack injects malicious client-side scripts (usually JavaScript) into trusted websites, which then execute in a victim's browser)
		w.Header().Set("X-Content-Type-Options", "nosniff")                                       // prevents browsers from mime sniffing (MIME sniffing (or content sniffing) is when a web browser inspects the actual bytes of a file to guess its MIME type (like image/jpeg, text/html) instead of relying solely on the Content-Type header provided by the server)
		w.Header().Set("Strict Transport Security", "max-age=63072000;includeSubDomains;preload") // enforces https for max age in seconds (will refuse to connect to site without https)
		w.Header().Set("Content-Security-Policy", "default-src 'self'")                           // allowing only resources from same resource
		w.Header().Set("Referrer-Policy", "no-referrer")                                          // controls how much referrer req is sent with request
		w.Header().Set("X-Powered-By", "Django")

		// more headers
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")

		next.ServeHTTP(w, r)
	})
}

/*
1️⃣ X-DNS-Prefetch-Control: off

	What it does
		Controls DNS prefetching by the browser.
		DNS prefetching resolves domain names before the user clicks a link.

	Why disable it?
		Prevents:
		Information leakage (browser resolving domains you might not visit)
		Unnecessary DNS queries

	Notes
		Mostly relevant for older browsers
		Modern browsers largely manage this automatically

2️⃣ X-Frame-Options: DENY

	What it’s supposed to do
		Prevents clickjacking
		Controls whether your site can be embedded in an <iframe>

	Valid values
		Value	Meaning
		DENY	Cannot be framed anywhere
		SAMEORIGIN	Can be framed by same origin only

	or (recommended via CSP):
		Content-Security-Policy: frame-ancestors 'none'

3️⃣ X-XSS-Protection: 1; mode=block ⚠️ (Mostly obsolete)

	What it does
		Enables browser’s built-in XSS filter
		Blocks rendering if XSS is detected

	Status
		❌ Deprecated
		Ignored by modern Chrome, Edge, Firefox
		Can cause security issues in older browsers
		Modern replacement
		Use Content Security Policy (CSP) instead

4️⃣ X-Content-Type-Options: nosniff ✅

	What it does
		Prevents MIME-type sniffing
		Forces browser to respect Content-Type
		Prevents
		Attacks where a file is served as: text/plain


	but executed as:
	application/javascript

5️⃣ Strict-Transport-Security (HSTS) ✅
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload

	What it does
		Forces HTTPS only
		Browser will never try HTTP for this domain
		Breakdown
		Directive	Meaning
		max-age=63072000	Enforce HTTPS for 2 years
		includeSubDomains	Applies to all subdomains
		preload	Allows inclusion in browser HSTS preload list

	⚠️ Important warning
	Once enabled:

		You cannot serve HTTP anymore
		Misconfiguration can brick your domain
		Only use if:
		HTTPS works perfectly everywhere

6️⃣ Content-Security-Policy (CSP) ✅🔥
Content-Security-Policy: default-src 'self'

	What it does
		Most powerful security header
		Controls where resources can load from
		This policy means

	Only allow:
		JS
		CSS
		Images
		Fonts
		Frames
		from the same origin
	Prevents
		XSS
		Data injection
		Malicious third-party scripts

	Real-world CSP is usually more detailed:
	Content-Security-Policy:
	default-src 'self';
	script-src 'self' https://trusted.cdn.com;
	object-src 'none';

7️⃣ Referrer-Policy
Referrer-Policy: no-referrer

	What it does
		Controls how much referrer information is sent
		Prevents leaking URLs, tokens, or sensitive paths

	Common values
		Value							Behavior
		no-referrer						Never send referrer
		strict-origin-when-cross-origin	Modern default
		same-origin						Only for same site

*/

// ------BASIC MIDDLEWARE SKELETON-----
/*
func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Do something BEFORE the handler
		fmt.Println("Before handler")

		// 2. Call the next handler
		next.ServeHTTP(w, r)

		// 3. Do something AFTER the handler
		fmt.Println("After handler")
	})
}
*/
