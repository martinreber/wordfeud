export interface LetterPlaySet {
    letter: string;
    count: number;
    value: number;
}

export interface Session {
    user: string;
    last_move_timestamp: string;
    session_start_timestamp: string;
    reminding_letters: number;
}

export interface SessionResponse {
    sessions: Session[];
}

export interface PlayedMove {
    letters: string;
    word: string;
    played_by_myself: boolean;
    timestamp: string;
}

export interface UserSession {
    letters_play_set: LetterPlaySet[];
    last_move_timestamp: string;
    session_start_timestamp: string;
    letter_overall_value: number;
    played_moves: PlayedMove[];
}

export interface WordCount {
    word: string;
    count: number;
}

export type WordCounts = WordCount[];