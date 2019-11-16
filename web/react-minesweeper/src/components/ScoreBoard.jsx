import React, {useState, useEffect} from 'react';

const ScoreBoard = React.memo((props) => {
    const [timer, setTimer] = useState(0);
    const [isActive, setActive] = useState(false);
    const toggle = () => {
        setActive(prevState => !prevState);
    }

    const reset = () => {
        setActive(false);
        setTimer(0);
    }

    if (props.start) {
        if (!isActive) {
            toggle();
        }
    } else if (props.stop) {
        if (isActive) {
            toggle();
        }
    } else if (props.reset) {
        if (timer !== 0) {
            reset();
        }
    }

    useEffect(() => {
        let myTimer;
        if (isActive) {
            myTimer = setInterval(() => setTimer(prevState => prevState + 1), 1000);
        } else if (!isActive && timer !== 0) {
            clearInterval(myTimer);
        }
        return () => clearInterval(myTimer);
    }, [isActive, timer]);

    return (
        <div className='scoreboard'>
            {((props.start || props.stop) && <div className='mines'>
            Mines: {props.numMines}
            </div>)}
            {((props.start || props.stop) && <Timer time={timer} />)}
        </div>
    );
});

const Timer = (props) => {
    return (<div className='timer'>
            Time: {props.time}
            </div>);
};

export default ScoreBoard;
