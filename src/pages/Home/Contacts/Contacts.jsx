import { useState } from "react";
import styles from "./Contacts.module.scss";
import { MessageCircle } from "lucide-react"; // иконка (lucide-react)

const faqs = [
    {
        question: "Как мога да се свържа с вас?",
        answer: "Можете да ни пишете чрез формата за контакт или директно на email адреса ни."
    },
    {
        question: "Предлагате ли поддръжка?",
        answer: "Да, предоставяме 24/7 поддръжка чрез имейл и чат."
    },
    {
        question: "Колко време отнема отговорът?",
        answer: "Обикновено отговаряме в рамките на 1-2 работни дни."
    },
    {
        question: "Колко време отнема отговорът?",
        answer: "Обикновено отговаряме в рамките на 1-2 работни дни."
    },
    {
        question: "Колко време отнема отговорът?",
        answer: "Обикновено отговаряме в рамките на 1-2 работни дни."
    },
    {
        question: "Колко време отнема отговорът?",
        answer: "Обикновено отговаряме в рамките на 1-2 работни дни."
    }
];

export default function Section4() {
    const [openIndex, setOpenIndex] = useState(null);

    const toggleFAQ = (index) => {
        setOpenIndex(openIndex === index ? null : index);
    };

    return (
        <div className={styles.wrapper}>
            {/* FAQ (ляво) */}
            <div className={styles.faq}>
                {faqs.map((faq, index) => (
                    <div
                        key={index}
                        className={`${styles.faq_item} ${openIndex === index ? styles.open : ""}`}
                        onClick={() => toggleFAQ(index)}
                    >
                        <div className={styles.question}>
                            <span>{faq.question}</span>
                            <span>{openIndex === index ? "−" : "+"}</span>
                        </div>
                        {openIndex === index && (
                            <div className={styles.answer}>{faq.answer}</div>
                        )}
                    </div>
                ))}
            </div>

            {/* Картичка с въпросите (дясно) */}
            <div className={styles.contact_box}>
                <MessageCircle size={60} className={styles.icon} />
                <h3>Въпроси?</h3>
                <p>Не се колебайте да ни пишете. Нашият екип ще ви отговори възможно най-бързо.</p>
                <button className={styles.ask_btn}>Попитай</button>
            </div>
        </div>
    );
}
