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
        Ok(_) => {}
        Err(e) => println!("error parsing header: {e:?}"),
    }

    // web_sys::console::log_1(&"XKeys from embedded wasm!".into());
    let pair = KeyPair::new_from_raw(nkeys::KeyPairType::Curve, buf);
    match pair {
        Ok(pair) => {
            let public_key = pair.public_key();
            let seed = pair.seed().unwrap();

            let xkey = json!({
                "public": public_key,
                "seed": seed,
            });

            web_sys::console::log_1(&public_key.into());
            web_sys::console::log_1(&seed.into());
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
    let mine = "SXAJ7S4BJAFGK6JYNEXSY7LFLF4LX57UVCOHFNOYDK3TMRIPFGTNFDPSKQ";
    let theirs = "XDPYKH2N6JV3TUWWBWM2HWJX4FYEHDUYMJCPSYDJRAURXSBVVZTV2BAG";
    let a = Uint8Array::new(&JsValue::NULL);
    //Uint8Array::copy_from(&a, &[0]);
    a.copy_from(&[
        120, 107, 118, 49, 50, 180, 249, 241, 115, 221, 5, 5, 106, 143, 46, 162, 228, 135, 89, 169,
        56, 77, 107, 136, 91, 253, 64, 31, 242, 36, 163, 146, 77, 7, 176, 42, 241, 236, 136, 117,
        179, 34, 128, 78, 74, 63, 48, 144, 207, 248, 240, 34, 154, 239, 103, 116, 137, 103, 1, 151,
        127, 52, 175, 163, 151, 178,
    ]);

    //let data = [120, 107, 118, 49, 50, 180, 249, 241, 115, 221, 5, 5, 106, 143, 46, 162, 228, 135, 89, 169, 56, 77, 107, 136, 91, 253, 64, 31, 242, 36, 163, 146, 77, 7, 176, 42, 241, 236, 136, 117, 179, 34, 128, 78, 74, 63, 48, 144, 207, 248, 240, 34, 154, 239, 103, 116, 137, 103, 1, 151, 127, 52, 175, 163, 151, 178]
    let b = decrypt_chunk(a, mine.to_string(), theirs.to_string());
    //let result = 2 + 2;
    assert_eq!(b, "");
}
