// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

export function GetAllSettings(): Promise<any> & { cancel(): void } {
    let $resultPromise = $Call.ByID(1761571667) as any;
    return $resultPromise;
}

export function SetSetting(key: string, value: string): Promise<void> & { cancel(): void } {
    let $resultPromise = $Call.ByID(3519942637, key, value) as any;
    return $resultPromise;
}
