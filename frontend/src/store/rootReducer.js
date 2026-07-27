import { combineReducers } from '@reduxjs/toolkit';
import chatsReducer from './slices/chatSlice';
import meSliceReducer from './slices/meSlice';
import chatDialogReducer from './slices/chatDialogSlice';

const appReducer = combineReducers({
  chats: chatsReducer,
  user: meSliceReducer,
  chatDialog: chatDialogReducer,
});

const rootReducer = (state, action) => {
  if (action.type === 'LOGOUT') {
    state = undefined;
  }

  return appReducer(state, action);
};

export default rootReducer;
