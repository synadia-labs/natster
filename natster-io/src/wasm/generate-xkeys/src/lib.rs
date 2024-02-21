mod utils;

extern crate web_sys;

use js_sys::Uint8Array;
use nkeys::KeyPair;
use serde_json::json;

use wasm_bindgen::prelude::*;
use wasm_bindgen_test::*;

#[wasm_bindgen]
pub fn get_xkeys() -> JsValue {
    let mut buf = [0u8; 32];
    let r = getrandom::getrandom(&mut buf);
    match r {
        Ok(_) => {
            web_sys::console::log_1(&format!("random: {:?}", buf.to_vec()).into());
        }
        Err(e) => println!("error parsing header: {e:?}"),
    }

    let pair = KeyPair::new_from_raw(nkeys::KeyPairType::Curve, buf);
    match pair {
        Ok(pair) => {
            // let public_key = pair.public_key();
            // let seed = pair.seed().unwrap();
            let public_key = "XDEUIW5UQGWQLYWM2GDGWJ5J2XPDBNA63B3PEWVRBRSXWVVSVBNGX6QT";
            let seed = "SXAPFTI242W2ZLCJQAFHTQ3T3VZKSJZZGKNZOTJV4WG7T2MYWYYW7N7OEY";
            web_sys::console::log_1(&format!("public key: {}", public_key).into());
            web_sys::console::log_1(&format!("seed: {}", seed).into());

            let xkey = json!({
                "public": public_key,
                "seed": seed,
            });

            JsValue::from_str(xkey.to_string().as_str())
        }
        Err(e) => JsValue::from_str(format!("Error parsing header: {}", e).as_str()),
    }
}

#[wasm_bindgen]
pub fn decrypt_chunk(encrypted_in: Uint8Array, seed: String, sender: String) -> JsValue {
    let input: Vec<u8> = encrypted_in.to_vec();
    let my_xkey = nkeys::XKey::from_seed(seed.as_str());
    let their_xkey = nkeys::XKey::from_public_key(sender.as_str());
    web_sys::console::log_1(&format!("mine: {}", seed).into());
    web_sys::console::log_1(&format!("theirs: {}", sender).into());
    web_sys::console::log_1(&format!("data: {:?}", input.as_slice()).into());

    match my_xkey {
        Ok(mine) => match their_xkey {
            Ok(theirs) => {
                let txk = mine.public_key();
                web_sys::console::log_1(&format!("mine_pub: {}", txk).into());
                let dec = mine.open(input.as_slice(), &theirs);
                match dec {
                    Ok(decrypted_chunk) => {
                        JsValue::from_str(String::from_utf8(decrypted_chunk).unwrap().as_str())
                    }
                    Err(e) => JsValue::from_str(format!("Error decrypting chunk: {}", e).as_str()),
                }
            }
            Err(e) => JsValue::from_str(format!("Error parsing their key: {}", e).as_str()),
        },
        Err(e) => JsValue::from_str(format!("Error parsing seed: {}", e).as_str()),
    }
}

#[wasm_bindgen_test]
fn it_works() {
    let mine = "SXAPFTI242W2ZLCJQAFHTQ3T3VZKSJZZGKNZOTJV4WG7T2MYWYYW7N7OEY";
    let theirs = "XACKDJCVN7UF7WTXNHR5R224FGVN56AL45UF2JC37FAOFYBY6WDFXIDR";
    let a = Uint8Array::new(&JsValue::NULL);
    a.copy_from(&[
        120, 107, 118, 49, 18, 112, 233, 144, 97, 92, 195, 225, 63, 111, 79, 221, 156, 131, 166,
        142, 122, 163, 45, 191, 47, 140, 20, 46, 189, 84, 191, 226, 224, 94, 239, 169, 223, 22,
        225, 32, 99, 241, 171, 249, 64, 132, 33, 225, 28, 225, 13, 237, 146, 183, 149, 117, 175,
    ]);

    let b = decrypt_chunk(a, mine.to_string(), theirs.to_string());
    console_log!("{:?}", b);
    assert_eq!(b, "hello kevin");
}
