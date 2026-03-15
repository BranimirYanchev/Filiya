import React from 'react';
import styles from './Home.module.scss';
import Carousel from './Carousel/Carousel';
import Navbar from '../../components/Navbar/Navbar';
import Contacts from './Contacts/Contacts';
import Footer from '../../components/Footer/Footer';

export default function Home({ onNavigate }) {
  return (
    <div className={styles.home}>
      <Navbar currentPage="home" onNavigate={onNavigate} />
      <section className={styles.section_1}>
        <div className={styles.content}>
          <h1>φιλία</h1>
          <p>Добре дошли във социалната мрежа</p>
        </div>
      </section>
      <section className={styles.section_2}>
        <div className={styles.content}>
          <h1>Най-популярни постове</h1>
          <Carousel />
        </div>
      </section>
      <section className={styles.section_3}>
        <div className={styles.content}>
          <h1>Информация</h1>
          <div className={styles.info_box}>
            <img src={require("./Carousel/images/img1.png")} className={styles.info_image} alt="Информация 1" />
            <div className={styles.info_content}>
              <h2>Заглавие 1</h2>
              <p>Това е място, където ще откриете вдъхновяващи истории, полезни съвети и красиви моменти, споделени от хора по целия свят. Потопете се в разнообразни теми – от пътешествия и кулинария до личностно развитие и изкуство. Всяка публикация носи емоция и ново знание.</p>
            </div>
          </div>
          <div className={`${styles.info_box} ${styles.reverse}`}>
            <img src={require("./Carousel/images/img2.png")} className={styles.info_image} alt="Информация 2" />
            <div className={styles.info_content}>
              <h2>Заглавие 2</h2>
              <p>Това е място, където ще откриете вдъхновяващи истории, полезни съвети и красиви моменти, споделени от хора по целия свят. Потопете се в разнообразни теми – от пътешествия и кулинария до личностно развитие и изкуство. Всяка публикация носи емоция и ново знание.</p>
            </div>
          </div>
          <div className={styles.info_box}>
            <img src={require("./Carousel/images/img3.png")} className={styles.info_image} alt="Информация 3" />
            <div className={styles.info_content}>
              <h2>Заглавие 3</h2>
              <p>Това е място, където ще откриете вдъхновяващи истории, полезни съвети и красиви моменти, споделени от хора по целия свят. Потопете се в разнообразни теми – от пътешествия и кулинария до личностно развитие и изкуство. Всяка публикация носи емоция и ново знание.</p>
            </div>
          </div>
        </div>
      </section>
      <section className={styles.section_4}>
        <div className={styles.content}>
          <h1>Контакти</h1>
          <Contacts/>
        </div>
      </section>
      <Footer />
    </div>
  );
}
