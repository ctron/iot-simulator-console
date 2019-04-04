import React from "react";

import {
    BackgroundImage, BackgroundImageSrc,
    Page, PageHeader, PageSection,
    DataList, DataListItem, DataListCell,
    Title, Brand, Alert,
} from "@patternfly/react-core";

import {
    ChartBar, ChartLabel, ChartLegend, ChartTheme, Chart, ChartGroup, ChartPie,
    ChartContainer,
} from "@patternfly/react-charts";

import {
    CubesIcon, AngleDoubleDownIcon, AngleDoubleUpIcon, ClockIcon, ExclamationTriangleIcon,
    CheckCircleIcon
} from '@patternfly/react-icons';

import "./App.css"
import brandImg from "./iot-simulator.svg"

import backgroundLg from "./assets/images/background_1200.jpg"
import backgroundSm from "./assets/images/background_768.jpg"
import backgroundXs from "./assets/images/background_576.jpg"
import backgroundSm2x from "./assets/images/background_768@2x.jpg"
import backgroundXs2x from "./assets/images/background_576@2x.jpg"

import backgroundFilter from "./assets/images/background-filter.svg"

class App extends React.Component {

    render() {
        const background = {
            [BackgroundImageSrc.lg]: backgroundLg,
            [BackgroundImageSrc.sm]: backgroundSm,
            [BackgroundImageSrc.sm2x]: backgroundSm2x,
            [BackgroundImageSrc.xs]: backgroundXs,
            [BackgroundImageSrc.xs2x]: backgroundXs2x,
            [BackgroundImageSrc.filter]: backgroundFilter + "#image_overlay",
        }
        const Header = (<PageHeader
            logo={<Brand alt="IoT Simulator" src={brandImg}/>}
        />)

        return (
            <React.Fragment>
                <BackgroundImage src={background}/>
                <Page header={Header}>
                    <PageSection>
                        <Home/>
                    </PageSection>
                </Page>
            </React.Fragment>
        )
    }
}

export default App;

class Home extends React.Component {

    constructor(props) {
        super(props);
        this.api = "//" + window.location.host + "/api";
        this.state = {
            overview: {
                tenants: []
            },
        };
        this.refreshData();
    }

    refreshData() {
        fetch(this.api + "/overview")
            .then(result => {
                return result.json()
            })
            .then(data => {
                console.log(data);
                this.setState({overview: data});
            })
    }

    componentDidMount() {
        this.interval = setInterval(() => this.refreshData(), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    static fixedNumber(n, fractionDigits) {
        if (n == null || isNaN(n)) {
            return "␀"
        }
        return n.toFixed(fractionDigits)
    }

    renderConsumers(tenant) {
        if (tenant.consumers == null) {
            return
        }

        const o = this;

        return tenant.consumers.map(function (consumer, i) {
            return (
                <DataListItem>
                    <DataListCell>
                        {o.renderGood(consumer)}&nbsp;
                        {consumer.type}
                    </DataListCell>
                    <DataListCell>
                        <CubesIcon/>&nbsp;<strong>{consumer.replicas} Pods</strong>
                    </DataListCell>
                    <DataListCell>
                        <AngleDoubleDownIcon/>&nbsp;
                        <strong title="msgs/s" data-toggle="tooltip" data-placement="top">
                            {o.renderSingleValue(consumer.messagesPerSecond, "msgs/s")}
                        </strong>
                    </DataListCell>
                    <DataListCell>&nbsp;</DataListCell>
                    <DataListCell width={2}>&nbsp;</DataListCell>
                </DataListItem>
            )
        })
    }

    renderErrorChart(producer) {
        return (<div className="chart-inline">
            <div>
                <ChartPie
                    animate={{duration: 500}}
                    containerComponent={<ChartContainer responsive={false}/>}
                    labels={datum => `${datum.x}: ${Home.fixedNumber(datum.y, 0)}`}
                    height={80} width={80}
                    padding={10}
                    data={producer.chartData}
                />
            </div>
            <ChartLegend
                orientation={"vertical"}
                data={producer.chartLegend}
                rowGutter={-8} gutter={20}
                itemsPerRow={2}
                height={80} width={200}
            />
        </div>)
    }

    renderConnectionChart(producer) {
        return (<div className="chart-inline">
            <div>
                <ChartPie
                    animate={{duration: 500}}
                    containerComponent={<ChartContainer responsive={false}/>}
                    labels={datum => `${datum.x}: ${Home.fixedNumber(datum.y, 0)}`}
                    height={80} width={80}
                    padding={10}
                    data={[
                        {"x": "connected", "y": producer.connectionsEstablished},
                        {"x": "disconnected", "y": producer.connectionsConfigured - producer.connectionsEstablished},
                    ]}
                />
            </div>
            <ChartLegend
                orientation={"vertical"}
                data={[{"name": "connected"}, {"name": "disconnected"}]}
                rowGutter={-8} gutter={20}
                itemsPerRow={2}
                height={80} width={200}
            />
        </div>)
    }

    renderSingleValue(value, tooltip) {
        return (<span
            title={tooltip} data-toggle="tooltip"
            data-placement="top">
            {(value != null) ? value.toFixed(0) : "␀"}
        </span>)
    }

    renderGood(common) {
        return (
            <div className="state-indicator">{common.good ? <CheckCircleIcon/> : <ExclamationTriangleIcon/>}</div>
        )
    }

    renderProducers(tenant) {

        const o = this

        if (tenant.producers == null) {
            return
        }
        return tenant.producers.map(function (producer, i) {
            return (
                <DataListItem className="chart-list">
                    <DataListCell>
                        {o.renderGood(producer)}&nbsp;
                        {producer.type + " / " + producer.protocol}
                    </DataListCell>
                    <DataListCell>
                        <CubesIcon/>&nbsp;
                        <strong>{producer.replicas} Pods</strong>
                    </DataListCell>
                    <DataListCell>
                        <AngleDoubleUpIcon/>&nbsp;
                        <strong>
                            {o.renderSingleValue(producer.messagesPerSecondSent, "msgs/s")}
                        </strong>
                    </DataListCell>
                    <DataListCell>
                        <ClockIcon/>&nbsp;
                        <strong>
                            {o.renderSingleValue(producer.messagesPerSecondConfigured, "msgs/s configured")}&nbsp;→&nbsp;
                            {o.renderSingleValue(producer.messagesPerSecondScheduled, "msgs/s scheduled")}
                        </strong>
                    </DataListCell>
                    {producer.protocol === "mqtt" ?
                        <DataListCell className="chart-cell" width={2}>
                            {o.renderConnectionChart(producer)}
                        </DataListCell>
                        :
                        <DataListCell className="chart-cell" width={2}>
                            {producer.chartData != null && producer.chartLegend != null ? o.renderErrorChart(producer) : ""}
                        </DataListCell>
                    }
                </DataListItem>
            );
        })
    }

    render() {
        const o = this;
        return (
            <div>
                {
                    this.state.overview.tenants.map(function (tenant, i) {
                        return (
                            <div>
                                <Title size="3xl">{tenant.name}</Title>
                                <DataList>
                                    {o.renderConsumers(tenant)}
                                    {o.renderProducers(tenant)}
                                </DataList>
                            </div>
                        )
                    })
                }
            </div>
        );
    }
}
