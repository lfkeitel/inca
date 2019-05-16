import * as $ from "jquery";

export interface IStatusResult {
    running: boolean
    totalDevices: number
    finished: number
    stage: Stage
}

export enum Stage {
    Default = "",
    PreScript = "pre-script",
    LoadingConfig = "loading-configuration",
    Grabbing = "grabbing",
    PostScript = "post-script",
}

export interface IPerformRunResult {
    status: string
    running: boolean
}

export interface IDeviceListResult {
    devices: Device[]
}

export interface Device {
    path: string
    name: string
    address: string
    proto: string
    conf_text: string
    manufacturer: string
}

export interface ErrorLine {
    etype: string
    time: string
    message: string
}

export interface ISaveConfigResult {
    success: boolean
    error?: string
}

export function checkStatus(callback: (data: IStatusResult) => void) {
    $.get('/api/status', {}, null, 'json')
        .done((data: IStatusResult) => callback(data));
}

export function performRun(callback: (data: IPerformRunResult) => void) {
    $.get('/api/runnow', {}, null, 'json')
        .done((data: IPerformRunResult) => callback(data));
}

export function getDeviceList(callback: (data: IDeviceListResult) => void) {
    $.get('/api/devicelist', {}, null, 'json')
        .done((data: IDeviceListResult) => callback(data));
}

export function getErrorLog(callback: (data: ErrorLine[]) => void) {
    $.get('/api/errorlog', { limit: 10 }, null, 'json')
        .done((data: ErrorLine[]) => callback(data));
}

export function runSingleDeviceGrab(
    address: string,
    brand: string,
    proto: string,
    name: string,
    callback: (data: IPerformRunResult) => void) {
    $.get('/api/singlerun', { hostname: address, name: name, proto: proto, brand: brand }, null, 'json')
        .done((data: IPerformRunResult) => callback(data));
}

export function saveDeviceList(listText: string, callback: (data: ISaveConfigResult) => void) {
    $.post('/api/savedevicelist', { text: encodeURIComponent(listText) }, null, "json")
        .done((data: ISaveConfigResult) => callback(data));
}

export function saveDeviceTypes(listText: string, callback: (data: ISaveConfigResult) => void) {
    $.post('/api/savedevicetypes', { text: encodeURIComponent(listText) }, null, "json")
        .done((data: ISaveConfigResult) => callback(data));
}
