// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {database} from '../models';
import {service} from '../models';

export function AddCard(arg1:number,arg2:string,arg3:number):Promise<void>;

export function Cleanup():Promise<void>;

export function DeleteCard(arg1:number,arg2:number):Promise<void>;

export function GetActiveTimeEntry(arg1:number,arg2:number):Promise<database.TimeEntry>;

export function GetAll(arg1:number,arg2:service.CardStatus):Promise<Array<database.ListCardsRow>>;

export function GetCardById(arg1:number,arg2:number):Promise<database.GetCardRow>;

export function StartCard(arg1:number,arg2:number):Promise<void>;

export function StopCard(arg1:number,arg2:number):Promise<void>;

export function UpdateCard(arg1:number,arg2:number,arg3:service.UpdateCardParams):Promise<void>;

export function UpdateCardStatus(arg1:number,arg2:number,arg3:service.CardStatus):Promise<void>;
