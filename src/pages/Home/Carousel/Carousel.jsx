import React, { useState, useRef, useEffect } from "react";
import "./Carousel.scss";
import leftArrow from './images/left-arr.svg';
import rightArrow from './images/right-arr.svg';

const images = [
    require("./images/img1.png"),
    require("./images/img2.png"),
    require("./images/img3.png"),
];

const Carousel = () => {
    const [current, setCurrent] = useState(0);
    const trackRef = useRef(null);
    const [height, setHeight] = useState(0);
    const [isAnimating, setIsAnimating] = useState(false);
    const animationTimeoutRef = useRef(null);

    const startAnimation = () => {
        setIsAnimating(true);

        if (animationTimeoutRef.current) {
            clearTimeout(animationTimeoutRef.current);
        }

        animationTimeoutRef.current = setTimeout(() => {
            setIsAnimating(false);
        }, 520);
    };

    const goToSlide = (targetIndex) => {
        if (isAnimating || targetIndex === current) return;
        setCurrent(targetIndex);
        startAnimation();
    };

    const prevSlide = () => {
        if (isAnimating) return;
        setCurrent((prev) => (prev === 0 ? images.length - 1 : prev - 1));
        startAnimation();
    };

    const nextSlide = () => {
        if (isAnimating) return;
        setCurrent((prev) => (prev === images.length - 1 ? 0 : prev + 1));
        startAnimation();
    };

    useEffect(() => {
        if (trackRef.current) {
            const activeSlide = trackRef.current.querySelector(".activeSlide img");

            if (activeSlide) {
                const observer = new ResizeObserver((entries) => {
                    for (let entry of entries) {
                        setHeight(entry.contentRect.height);
                    }
                });

                observer.observe(activeSlide);

                // вземи височината и веднага
                setHeight(activeSlide.offsetHeight);

                return () => observer.disconnect();
            }
        }
    }, [current]);

    useEffect(() => {
        return () => {
            if (animationTimeoutRef.current) {
                clearTimeout(animationTimeoutRef.current);
            }
        };
    }, []);


    return (
        <div className="carousel">
            <div
                ref={trackRef}
                className="carousel__track"
                style={{ height: `${height}px`, transition: "height 0.3s ease" }}
            >
                {images.map((img, index) => {
                    let position = "nextSlide";
                    if (index === current) position = "activeSlide";
                    if (index === current - 1 || (current === 0 && index === images.length - 1))
                        position = "prevSlide";

                    return (
                        <div key={index} className={`carousel__slide ${position}`}>
                            <img src={img} alt={`Пост ${index + 1}`} loading="lazy" />
                        </div>
                    );
                })}
            </div>

            <div className="carousel__dots">
                <button className="carousel__arrow left" onClick={prevSlide} disabled={isAnimating} aria-label="Предишен слайд">
                    <img src={leftArrow} alt="prev" />
                </button>

                {images.map((_, idx) => (
                    <button
                        key={idx}
                        className={`dot ${idx === current ? "active" : ""}`}
                        onClick={() => goToSlide(idx)}
                        aria-label={`Отиди на слайд ${idx + 1}`}
                    />
                ))}

                <button className="carousel__arrow right" onClick={nextSlide} disabled={isAnimating} aria-label="Следващ слайд">
                    <img src={rightArrow} alt="next" />
                </button>
            </div>
        </div>
    );
};

export default Carousel;
