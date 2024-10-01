import { actions } from "./index";

export const reducer = (state, action) => {
    switch (action.type) {
        case actions.SHOW_OPERATION_TOAST:
            return {
                ...state,
                toastMessage: action.data
            };
        default:
            return state;
    }
};
