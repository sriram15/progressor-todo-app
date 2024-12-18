// This file is intended to be an abstraction for the frontend to not know about @wailsjs APIs
// In the future, this will help in converting the frontend code to use fetch() request instead of relying on API
// Right now, we are just going to export the same as wails App
import {
  GetAll,
  GetOpenCards,
  GetCardById,
  AddCard,
  DeleteCard,
  UpdateCard,
  UpdateCardStatus,
  StartCard,
  StopCard,
} from "@wailsjs/go/service/cardService";

export {
  GetAll,
  GetOpenCards,
  GetCardById,
  AddCard,
  DeleteCard,
  UpdateCard,
  UpdateCardStatus,
  StartCard,
  StopCard,
};
