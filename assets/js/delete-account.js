"use strict";
/*
 The account deletion page (/location/delete-account).

 Everything happens in the browser on purpose: unwrapping the user's secret key
 requires their password, and their password must never reach a Zood server.
 We talk to oscar directly with the shared client in zood-js.

 All of the copy lives in the page's markup, rendered by the l10n package. This
 file only shows and hides the pieces of it.
*/
var deleteAccount;
(function (deleteAccount) {
    const stepIDs = ["step-confirm", "step-login", "step-warning", "step-done"];
    // The rule oscar enforces when an account is created (validUsernamePattern
    // plus the 32 character limit in oscar's users.go), applied after the same
    // trimming and lowercasing. Every account that exists matches it.
    const validUsernamePattern = /^[a-z0-9]{5,32}$/;
    // The access token we get from a successful login. It only lives in memory,
    // and only long enough to delete the account.
    let accessToken = null;
    const sodiumScriptPath = "/js/zood-js/libs/sodium-sumo.0.7.15.js";
    /*
     zood.js reads window.sodium as soon as it's evaluated, long before
     libsodium.js replaces it with the real thing, so it needs a placeholder.
     This assignment has to happen before zood.js runs, which is why this script
     is included ahead of it.
    */
    window.sodium = {};
    let sodiumLoad = null;
    /*
     loadSodium adds libsodium.js to the page and resolves once its functions are
     safe to call. We add the script ourselves, instead of with a tag in the
     markup, so we find out when it fails: the page's CSP forbids inline onerror
     handlers, and libsodium.js only ever reports success, so without this a
     failed download leaves the visitor watching the spinner forever.

     A failed download is worth retrying on the next submit, so it clears
     sodiumLoad. A failure to initialize won't get better, so it doesn't.
    */
    function loadSodium() {
        if (sodiumLoad != null) {
            return sodiumLoad;
        }
        sodiumLoad = new Promise((resolve, reject) => {
            const script = document.createElement("script");
            script.src = sodiumScriptPath;
            script.addEventListener("load", function () {
                // Once it has run, libsodium.js has replaced window.sodium with an
                // object whose ready promise settles when initialization does. It
                // rejects if neither the wasm nor the asm.js build will start.
                const ready = window.sodium.ready;
                if (ready == null) {
                    reject(new Error("libsodium.js ran without defining sodium.ready"));
                    return;
                }
                ready.then(resolve, reject);
            });
            script.addEventListener("error", function () {
                script.remove();
                sodiumLoad = null;
                reject(zood.Error.networkError());
            });
            document.head.appendChild(script);
        });
        return sodiumLoad;
    }
    // UserFacingError carries a message that's already fit to show the user.
    class UserFacingError {
        constructor(message) {
            this.message = message;
        }
    }
    function element(id) {
        const el = document.getElementById(id);
        if (el == null) {
            throw new Error("Missing element '" + id + "'");
        }
        return el;
    }
    // message returns one of the error strings the server rendered into the
    // page's data attributes.
    function message(name) {
        const value = element("delete-account").dataset[name];
        return value == null ? "" : value;
    }
    function messageForError(err) {
        // The page only ever shows a handful of generic messages, so put the
        // real failure somewhere a developer can find it.
        console.error("Zood account deletion failed:", err);
        if (err instanceof UserFacingError) {
            return err.message;
        }
        if (err instanceof zood.Error) {
            switch (err.code) {
                case -2 /* zood.ErrorCode.network */:
                    return message("networkError");
                case 16 /* zood.ErrorCode.userNotFound */:
                case 11 /* zood.ErrorCode.notFound */:
                case 19 /* zood.ErrorCode.loginFailed */:
                    // Don't tell the visitor which half they got wrong.
                    return message("incorrectLogin");
            }
        }
        return message("unexpectedError");
    }
    function showStep(id) {
        for (const stepID of stepIDs) {
            element(stepID).hidden = stepID != id;
        }
    }
    function setProgressVisible(visible) {
        element("progress").hidden = !visible;
    }
    function showError(id, msg) {
        const el = element(id);
        el.textContent = msg;
        el.hidden = false;
    }
    function hideError(id) {
        const el = element(id);
        el.textContent = "";
        el.hidden = true;
    }
    function setLoginFormEnabled(enabled) {
        element("username").disabled = !enabled;
        element("password").disabled = !enabled;
        element("login-button").disabled = !enabled;
    }
    // nextFrame gives the browser a chance to paint before we hand the thread
    // to something slow.
    function nextFrame() {
        return new Promise((resolve) => {
            requestAnimationFrame(() => {
                setTimeout(resolve, 0);
            });
        });
    }
    /*
     login performs oscar's challenge/response login and returns an access
     token. It's the same sequence the Android app uses: fetch a challenge and
     the account's key derivation parameters, stretch the password to unwrap the
     account's secret key, then encrypt the challenge to the server's public key.
    */
    async function login(username, password) {
        const client = new zood.Client(null);
        const authChallenge = await client.getAuthenticationChallenge(username);
        const user = authChallenge.user;
        /*
         This is deliberately expensive - accounts created by the Android app
         use 30 passes over 64 MB - and it blocks this thread the whole time,
         which is why the caller shows a progress indicator first.
        */
        const passwordHash = zood.stretchPassword(sodium.crypto_secretbox_KEYBYTES, password, user.passwordSalt, user.passwordHashOperationsLimit, user.passwordHashMemoryLimit, user.passwordHashAlgorithm);
        let secretKey;
        try {
            secretKey = zood.symmetricKeyDecrypt(user.wrappedSecretKey, user.wrappedSecretKeyNonce, passwordHash);
        }
        catch (err) {
            // The wrapped secret key only fails to open if the password was wrong.
            throw new UserFacingError(message("incorrectLogin"));
        }
        const serverKey = await client.getServerPublicKey();
        const finishedChallenge = new zood.FinishedAuthenticationChallenge(zood.publicKeyEncrypt(authChallenge.challenge, serverKey.public_key, secretKey), zood.publicKeyEncrypt(authChallenge.creationDate, serverKey.public_key, secretKey));
        const loginResponse = await client.completeAuthenticationChallenge(username, finishedChallenge);
        return loginResponse.accessToken;
    }
    function onBeginClicked() {
        showStep("step-login");
        element("username").focus();
    }
    async function onLoginSubmitted(evt) {
        evt.preventDefault();
        hideError("login-error");
        const passwordField = element("password");
        // oscar lowercases usernames before looking them up, so a visitor whose
        // keyboard capitalized the first letter should still get in.
        const username = element("username").value.trim().toLowerCase();
        const password = passwordField.value;
        if (username == "" || password == "") {
            return;
        }
        // Catch typos here, before they turn into a confusing error from oscar
        // (or, for a '/' or '#', a request to the wrong URL), and before we spend
        // seconds stretching the password for an account that can't exist.
        if (!validUsernamePattern.test(username)) {
            showError("login-error", message("invalidUsername"));
            return;
        }
        setLoginFormEnabled(false);
        setProgressVisible(true);
        try {
            await loadSodium();
            await nextFrame();
            accessToken = await login(username, password);
            passwordField.value = "";
            showStep("step-warning");
        }
        catch (err) {
            showError("login-error", messageForError(err));
        }
        finally {
            setProgressVisible(false);
            setLoginFormEnabled(true);
        }
    }
    async function onDeleteClicked() {
        if (accessToken == null) {
            return;
        }
        hideError("delete-error");
        const button = element("delete-button");
        button.disabled = true;
        setProgressVisible(true);
        try {
            await new zood.Client(accessToken).deleteAccount();
            accessToken = null;
            showStep("step-done");
        }
        catch (err) {
            showError("delete-error", messageForError(err));
            button.disabled = false;
        }
        finally {
            setProgressVisible(false);
        }
    }
    document.addEventListener("DOMContentLoaded", function () {
        element("begin-button").addEventListener("click", onBeginClicked);
        element("login-form").addEventListener("submit", onLoginSubmitted);
        element("delete-button").addEventListener("click", onDeleteClicked);
        // Start the download now, so it has usually finished by the time the
        // visitor has typed their password. Any failure is reported when they
        // submit the login form, which calls loadSodium again.
        loadSodium().catch(function () { });
    });
})(deleteAccount || (deleteAccount = {}));
